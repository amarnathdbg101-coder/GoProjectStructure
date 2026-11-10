package model

import (
	"time"
)

// __PASCAL__ represents the database entity for __PLURAL__
type __PASCAL__ struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Create__PASCAL__Request defines input payload for creating a new __PASCAL__.
type Create__PASCAL__Request struct {
	Title       string `json:"title" validate:"required,min=2,max=150"`
	Description string `json:"description" validate:"max=1000"`
}

// Update__PASCAL__Request defines input payload for updating an existing __PASCAL__.
type Update__PASCAL__Request struct {
	Title       string `json:"title" validate:"omitempty,min=2,max=150"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}
