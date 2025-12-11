package entity

import "time"

type Note struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	OrganizationID string    `json:"organization_id"`
	UserID         string    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}
