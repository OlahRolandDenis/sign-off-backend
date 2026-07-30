package api

import (
	"Desktop/signoff/internal/auth"
	"Desktop/signoff/internal/db"
	"Desktop/signoff/internal/email"
	"Desktop/signoff/internal/models"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Test struct {
	Status   string
	Message  string
	AgencyID int64
}

type DataAgencyR struct {
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Plan       string    `json:"plan"`
	Created_at time.Time `json:"created_at"`
	AgencyID   int64     `json:"agency_id"`
	Api_keys   int       `json:"api_keys"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error bad method", http.StatusMethodNotAllowed)
		return
	}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error cant process request", http.StatusBadRequest)
		return
	}
	var re models.LoginRequest
	err = json.Unmarshal(res, &re)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	var agency *models.Agency

	agency, err = db.GetAgencyByEmail(r.Context(), re.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(agency.PasswordHash), []byte(re.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !agency.EmailVerified {
		http.Error(w, "Please verify your email first. Check your inbox.", http.StatusForbidden)
		return
	}

	token, err := auth.GenerateJWT(agency.ID)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var rez models.LoginResponse
	rez.Agency = *agency
	rez.Token = token
	data, err := json.Marshal(rez)
	if err != nil {
		http.Error(w, "Error reading body request", http.StatusBadRequest)
		return
	}

	w.Write(data)

}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error bad method", http.StatusMethodNotAllowed)
		return
	}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error cant process request", http.StatusBadRequest)
		return
	}

	var re models.RegisterRequest
	err = json.Unmarshal(res, &re)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if re.Email == "" || re.Name == "" || re.Password == "" {
		http.Error(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(re.Password), bcrypt.DefaultCost)

	id, err := db.CreateAgency(r.Context(), re.Email, re.Name, string(passHash))
	if err != nil {
		http.Error(w, "Error creating New Agency", http.StatusBadRequest)
		return
	}

	verifyToken := uuid.New().String()
	var query = `UPDATE agencies SET verification_token=$1 WHERE id=$2`
	_, err = db.Pool.Exec(r.Context(), query, verifyToken, id)
	if err == nil {
		go func() {
			_ = email.SendVerificationEmail(re.Email, verifyToken)
		}()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var t Test
	t.AgencyID = id
	t.Message = "Account created"
	t.Status = "success"

	data, err := json.Marshal(t)
	if err != nil {
		http.Error(w, "Couldnt make json", http.StatusBadRequest)
		return
	}

	w.Write(data)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error bad method", http.StatusMethodNotAllowed)
		return
	}

	agencyID, ok := r.Context().Value(AgencyIDKey).(int64)
	if ok == false {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rez, err := db.GetAgencyByID(r.Context(), agencyID)
	if err != nil {
		http.Error(w, "Agency not found", http.StatusNotFound)
		return
	}

	apikeys, err := db.GetAPIKeys(r.Context(), agencyID)
	if err != nil {
		http.Error(w, "Error getting apikeys", http.StatusBadRequest)
		return
	}

	var nr int
	nr = len(apikeys)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var t DataAgencyR
	t.AgencyID = rez.ID
	t.Api_keys = nr
	t.Created_at = rez.CreatedAt
	t.Email = rez.Email
	t.Name = rez.Name
	t.Plan = rez.Plan

	data, err := json.Marshal(t)
	if err != nil {
		http.Error(w, "Couldnt make json", http.StatusBadRequest)
		return
	}

	w.Write(data)

}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Error bad method", http.StatusMethodNotAllowed)
		return
	}

	agencyID, ok := r.Context().Value(AgencyIDKey).(int64)
	if ok == false {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := db.DeleteAgency(r.Context(), agencyID)
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing verification token", http.StatusBadRequest)
		return
	}
	var existingToken string

	var query = `Select verification_token FROM agencies WHERE verification_token=$1`
	err := db.Pool.QueryRow(r.Context(), query, token).Scan(&existingToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	query = `UPDATE agencies SET email_verified=true, verified_at = NOW(), verification_token = NULL WHERE verification_token = $1`
	_, err = db.Pool.Exec(r.Context(), query, token)
	if err != nil {
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("internal/templates/verify_succes.html")
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)

}
