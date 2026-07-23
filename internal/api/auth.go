package api

import (
	"Desktop/signoff/internal/auth"
	"Desktop/signoff/internal/db"
	"Desktop/signoff/internal/models"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Test struct {
	Status   string
	Message  string
	AgencyID int64
}

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
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

	token := auth.Store.CreateToken(agency.ID)

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
