package model

import "time"

// User represents the core user database entity.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone,omitempty"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Role         string    `json:"role"` // 'user', 'admin', etc.
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest holds input payload for user registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
}

// LoginRequest holds input credentials for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse returns token and user payload upon successful authentication.
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// UserResponse is a public safe view of a user profile.
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Role      string `json:"role"`
}
