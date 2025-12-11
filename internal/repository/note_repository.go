package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/entity"
)

type NoteRepository interface {
	Create(note *entity.Note) error
	GetAll(organizationID string) ([]entity.Note, error)
	GetByID(id int, organizationID string) (*entity.Note, error)
}

type postgresNoteRepository struct {
	db *sql.DB
}

func NewPostgresNoteRepository(db *sql.DB) NoteRepository {
	return &postgresNoteRepository{db: db}
}

func (r *postgresNoteRepository) Create(note *entity.Note) error {
	query := `
		INSERT INTO notes (title, content, organization_id, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(query, note.Title, note.Content, note.OrganizationID, note.UserID).
		Scan(&note.ID, &note.CreatedAt)

	if err != nil {
		return fmt.Errorf("gagal insert note: %w", err)
	}
	return nil
}

func (r *postgresNoteRepository) GetAll(organizationID string) ([]entity.Note, error) {
	query := `SELECT id, title, content, organization_id, user_id, created_at FROM notes WHERE organization_id = $1`

	rows, err := r.db.Query(query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []entity.Note
	for rows.Next() {
		var n entity.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.OrganizationID, &n.UserID, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *postgresNoteRepository) GetByID(id int, organizationID string) (*entity.Note, error) {
	query := `
		SELECT id, title, content, organization_id, user_id, created_at 
		FROM notes 
		WHERE id = $1 AND organization_id = $2
	`

	var n entity.Note
	err := r.db.QueryRow(query, id, organizationID).
		Scan(&n.ID, &n.Title, &n.Content, &n.OrganizationID, &n.UserID, &n.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}
