package models

import (
	"database/sql"
	"time"
)

type APIKey struct {
	ID        int64        `json:"id"`
	AgencyID  int64        `json:"agency_id"`
	Key       string       `json:"key"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	LastUsed  sql.NullTime `json:"last_used"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}
