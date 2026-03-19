package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string    `db:"name" json:"name" example:"John Doe"`
	Email     string    `db:"email" json:"email" example:"john.doe@example.com"`
	Password  string    `db:"password" json:"-" swaggerignore:"true"`
	Token     *string   `db:"token" json:"token,omitempty"`
	Active    bool      `db:"active" json:"active" example:"true"`
	CreatedAt time.Time `db:"created_at" json:"created_at" example:"2026-01-01T12:00:00Z"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" example:"2026-01-01T12:00:00Z"`
}
