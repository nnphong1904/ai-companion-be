package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents an app user.
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SignupRequest is the payload for user registration.
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is returned after successful authentication.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
