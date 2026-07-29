package api

import (
	"Desktop/signoff/internal/db"
	"Desktop/signoff/internal/email"
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type data struct {
	Title   string
	Content string
	Token   string
}

type data_status struct {
	Decision string
}

func ShowApprovePage(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Error extracting token", http.StatusBadRequest)
		return
	}

	var pool = db.Pool
	var ctx = context.Background()

	var query = `SELECT a.title, a.content, a.status, a.expires_at FROM approvals a WHERE a.token=$1`
	rez := pool.QueryRow(ctx, query, token)

	var title, status, content string
	var expires time.Time

	err := rez.Scan(&title, &content, &status, &expires)
	if err != nil {
		http.Error(w, "Approval not found", http.StatusNotFound)
		return
	}

	if status != "pending" {
		http.Error(w, "Approval Already Finished", http.StatusBadRequest)
		return
	}

	if status == "pending" && time.Now().After(expires) {
		http.Error(w, "Approval Expired", http.StatusGone)
		return
	}

	tmpl, err := template.ParseFiles("internal/templates/approve.html")
	if err != nil {
		http.Error(w, "Approval Expired", http.StatusInternalServerError)
		return
	}

	var de data
	de.Content = content
	de.Title = title
	de.Token = token

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	err = tmpl.Execute(w, de)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

}

func ProcessDecision(w http.ResponseWriter, r *http.Request) {

	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Error extracting token", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Error method invalid", http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var des = r.FormValue("decision")
	var com = r.FormValue("comment")
	var ctx = context.Background()
	var expiresAt time.Time

	if des != "approved" && des != "rejected" {
		http.Error(w, "Invalid decision value", http.StatusBadRequest)
		return
	}

	var q = `SELECT a.status, a.callback_url, a.expires_at, ag.email FROM approvals a JOIN agencies ag on ag.id=a.agency_id WHERE a.token=$1`
	var pool = db.Pool
	var status string
	var callbackurl string
	var agencyemail string

	rez := pool.QueryRow(ctx, q, token)
	err = rez.Scan(&status, &callbackurl, &expiresAt, &agencyemail)
	if err != nil {
		http.Error(w, "Approval not found", http.StatusNotFound)
		return
	}

	if status != "pending" {
		http.Error(w, "Approval Already Finished", http.StatusBadRequest)
		return
	}

	if time.Now().After(expiresAt) {
		http.Error(w, "Approval expired", http.StatusGone)
		return
	}

	q = `UPDATE approvals SET decision=$1, comment=$2, status=$3, decided_at=$4 WHERE token=$5`
	_, err = pool.Exec(ctx, q, des, com, des, time.Now(), token)
	if err != nil {
		http.Error(w, "Failed to save decision", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("internal/templates/confirm.html")
	if err != nil {
		http.Error(w, "Approval Expired", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	var tem_data data_status
	tem_data.Decision = des

	err = tmpl.Execute(w, tem_data)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

	if callbackurl != "" {
		go func() {
			err := SendCallback(callbackurl, token, des, com)
			if err != nil {
				log.Printf("Callback failed: %v", err)
				email.SendCallbackFailedEmail(agencyemail, token, des, err.Error())
			}
		}()
	}

}
