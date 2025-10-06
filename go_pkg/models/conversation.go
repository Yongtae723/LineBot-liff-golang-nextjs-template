package models

import "time"

type Conversation struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)
