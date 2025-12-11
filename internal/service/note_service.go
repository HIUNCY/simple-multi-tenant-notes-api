package service

import (
	"errors"
	"log"
	"time"

	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/entity"
	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/repository"
)

type NoteService interface {
	CreateNote(title, content, orgID, userID string) (*entity.Note, error)
	GetNotes(orgID string) ([]entity.Note, error)
	GetNoteByID(id int, orgID string) (*entity.Note, error)
}

type noteService struct {
	repo  repository.NoteRepository
	audit repository.AuditRepository
}

func NewNoteService(repo repository.NoteRepository, audit repository.AuditRepository) NoteService {
	return &noteService{repo: repo, audit: audit}
}

func (s *noteService) CreateNote(title, content, orgID, userID string) (*entity.Note, error) {
	if title == "" {
		return nil, errors.New("judul catatan tidak boleh kosong")
	}
	if orgID == "" || userID == "" {
		return nil, errors.New("identitas user/organisasi tidak valid")
	}

	newNote := &entity.Note{
		Title:          title,
		Content:        content,
		OrganizationID: orgID,
		UserID:         userID,
	}

	err := s.repo.Create(newNote)
	if err != nil {
		return nil, err
	}

	go func() {
		auditLog := &entity.AuditLog{
			Action:         "CREATE_NOTE",
			NoteID:         newNote.ID,
			OrganizationID: orgID,
			UserID:         userID,
			Details:        "Note created with title: " + title,
			Timestamp:      time.Now(),
		}

		if err := s.audit.Create(auditLog); err != nil {
			log.Printf("Gagal menulis Audit Log ke Mongo: %v", err)
		} else {
			log.Printf("Audit Log tersimpan di Mongo untuk Note ID: %d", newNote.ID)
		}
	}()

	return newNote, nil
}

func (s *noteService) GetNotes(orgID string) ([]entity.Note, error) {
	return s.repo.GetAll(orgID)
}

func (s *noteService) GetNoteByID(id int, orgID string) (*entity.Note, error) {
	return s.repo.GetByID(id, orgID)
}
