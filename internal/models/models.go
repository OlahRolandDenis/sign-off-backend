package models

import "time"

type Approval struct {
	ID            int64     `json:"id"`
	Token         string    `json:"token"`
	CallbackURL   string    `json:"callback_url"`
	ClientEmail   string    `json:"client_email"`
	AgencyID      int64     `json:"agency_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Status        string    `json:"status"`
	Decision      string    `json:"decision"`
	Comment       string    `json:"comment"`
	CreatedAt     time.Time `json:"created_at"`
	DecidedAt     time.Time `json:"decided_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	NudgeSent     bool      `json:"nudge_sent"`
	HangAlertSent bool      `json:"hang_alert_sent"`
}

type WebHookRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	CallbackURL string `json:"callback_url"`
	ClientEmail string `json:"client_email"`
}

type WebHookResponse struct {
	Status      string `json:"status"`
	Token       string `json:"token"`
	ApprovalURL string `json:"approval_url"`
	Message     string `json:"message"`
}
