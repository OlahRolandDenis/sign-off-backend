package api

import (
	"Desktop/signoff/internal/db"
	"Desktop/signoff/internal/email"
	"Desktop/signoff/internal/models"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error reading body request", http.StatusBadRequest)
		return
	}

	res, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body request", http.StatusBadRequest)
		return
	}

	agencyID, ok := r.Context().Value(AgencyIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, err := db.GetPlan(r.Context(), agencyID)
	if err != nil {
		http.Error(w, "Failed to verify plan", http.StatusInternalServerError)
		return
	}

	if plan == "free" {
		count, err := db.CountApprovalsInLast30Days(r.Context(), agencyID)
		if err != nil {
			http.Error(w, "Failed to count approvals", http.StatusInternalServerError)
			return
		}
		if count >= 50 {
			http.Error(w, "Free plan limit reached (50 approvals/month). Please upgrade to continue.", http.StatusPaymentRequired)
			return
		}
	}

	var response models.WebHookRequest

	err = json.Unmarshal(res, &response)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if response.Title == "" || response.Content == "" || response.ClientEmail == "" || response.CallbackURL == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	var token = uuid.New().String()

	var model models.Approval
	model.Token = token
	model.CallbackURL = response.CallbackURL
	model.ClientEmail = response.ClientEmail
	model.AgencyID = agencyID
	model.Title = response.Title
	model.Content = response.Content
	model.Status = "pending"
	model.CreatedAt = time.Now().UTC() //Fortare UTC
	model.ExpiresAt = time.Now().UTC().Add(48 * time.Hour)

	var query = `INSERT INTO approvals (token, callback_url, client_email, agency_id, title, content, status, created_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	var pool = db.Pool
	ctx := context.Background()

	_, err = pool.Exec(ctx, query, model.Token, model.CallbackURL, model.ClientEmail, model.AgencyID, model.Title, model.Content, model.Status, model.CreatedAt, model.ExpiresAt)
	if err != nil {
		http.Error(w, "Failed to save approval", http.StatusInternalServerError)
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	apvURL := baseURL + "/approve/" + token
	err = email.SendEmail(apvURL, response.Title, response.Content, response.ClientEmail)
	if err != nil {
		log.Printf("Email error: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	var rez models.WebHookResponse
	rez.Status = "success"
	rez.Token = token
	rez.ApprovalURL = baseURL + "/approve/" + token
	rez.Message = "Approval created successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data, err := json.Marshal(rez)
	if err != nil {
		http.Error(w, "Error reading body request", http.StatusBadRequest)
		return
	}
	w.Write(data)

}
