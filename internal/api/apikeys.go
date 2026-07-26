package api

import (
	"Desktop/signoff/internal/db"
	"Desktop/signoff/internal/models"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DataAPI struct {
	ID   int64
	Key  string
	Name string
}

type DataDel struct {
	Message string
}

func GenerateAPIKey() string {
	var b = make([]byte, 32)
	rand.Read(b)
	return "sk_live_" + base64.URLEncoding.EncodeToString(b)

}

func CreateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error bad method", http.StatusMethodNotAllowed)
		return
	}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error cant process request", http.StatusBadRequest)
		return
	}

	var re models.CreateAPIKeyRequest
	err = json.Unmarshal(res, &re)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if re.Name == "" {
		http.Error(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}

	//ok because returns an interface
	agencyID, ok := r.Context().Value(AgencyIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var key = GenerateAPIKey()
	id, err := db.CreateAPIKey(r.Context(), agencyID, re.Name, key)
	if err != nil {
		http.Error(w, "Failed to create API key", http.StatusInternalServerError)
		return
	}

	var rez DataAPI
	rez.ID = id
	rez.Key = key
	rez.Name = re.Name
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(rez)
	if err != nil {
		http.Error(w, "Error reading body request", http.StatusBadRequest)
		return
	}

	w.Write(data)
}

func ListAPIKeysHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error bad method", http.StatusMethodNotAllowed)
		return
	}

	//ok because returns an interface
	agencyID, ok := r.Context().Value(AgencyIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	keys, err := db.GetAPIKeys(r.Context(), agencyID)
	if err != nil {
		http.Error(w, "Failed to fetch API keys", http.StatusInternalServerError)
		return
	}

	for i := range keys {
		if len(keys[i].Key) > 12 {
			keys[i].Key = keys[i].Key[:12] + "..."
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(keys)
	if err != nil {
		http.Error(w, "Error reading body request", http.StatusBadRequest)
		return
	}

	w.Write(data)

}

func DeleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Error bad method", http.StatusMethodNotAllowed)
		return
	}

	agencyID, ok := r.Context().Value(AgencyIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing key ID", http.StatusBadRequest)
		return
	}

	newid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid key ID", http.StatusBadRequest)
		return
	}

	err = db.DeleteAPIKey(r.Context(), newid, agencyID)
	if err != nil {
		http.Error(w, "Failed to delete API key", http.StatusInternalServerError)
		return
	}

	var da DataDel
	da.Message = "API key revoked successfully"
	data, err := json.Marshal(da)
	if err != nil {
		http.Error(w, "Error reading body request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(data)
}
