package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	LineID    string    `json:"line_id" db:"line_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

