package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// User mirrors the auth-service user entity. quiz-app is read-only over HTTP;
// the source of truth is auth-service (Postgres). See CLAUDE.md §E#1.
type User struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EmailAddress string    `json:"email_address"`
	UserID       string    `json:"user_id"`
	Permission   []string  `json:"permission"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthClaims struct {
	UserID     string   `json:"user_id"`
	EmailID    string   `json:"email_id"`
	AuthorMail string   `json:"author_mail"`
	Exp        int64    `json:"exp"`
	Permission []string `json:"permission"`
	jwt.RegisteredClaims
}

// UpdateMeRequest — self-service update body forwarded to auth-service.
type UpdateMeRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Password  *string `json:"password,omitempty"`
}

// AdminUpdateUserRequest — admin update body forwarded to auth-service.
type AdminUpdateUserRequest struct {
	FirstName  *string   `json:"first_name,omitempty"`
	LastName   *string   `json:"last_name,omitempty"`
	Permission *[]string `json:"permission,omitempty"`
}

// PaginatedUsers — standard list response shape (CLAUDE.md G1 #10).
type PaginatedUsers struct {
	Items []User `json:"items"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
