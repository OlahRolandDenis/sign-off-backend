package models

import "time"

type Agency struct {
	ID                int64      `json:"id"`
	Email             string     `json:"email"`
	Name              string     `json:"name"`
	PasswordHash      string     `json:"-"`
	Plan              string     `json:"plan"`
	CreatedAt         time.Time  `json:"created_at"`
	EmailVerified     bool       `json:"email_verified"`
	VerificationToken *string    `json:"-"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Agency Agency `json:"agency"`
}
