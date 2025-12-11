package entity

import (
	"time"
)

type AuditLog struct {
	Action         string    `bson:"action"`
	NoteID         int       `bson:"note_id"`
	OrganizationID string    `bson:"organization_id"`
	UserID         string    `bson:"user_id"`
	Details        string    `bson:"details"`
	Timestamp      time.Time `bson:"timestamp"`
}
