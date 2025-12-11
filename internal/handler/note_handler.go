package handler

import (
	"net/http"
	"strconv"

	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/service"
	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	service service.NoteService
}

func NewNoteHandler(service service.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

type createNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

func (h *NoteHandler) Create(c *gin.Context) {
	var req createNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("userID")
	org, _ := c.Get("orgID")
	userIDStr := user.(string)
	orgIDStr := org.(string)

	if orgIDStr == "" || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User dan Organization ID tidak valid"})
		return
	}

	note, err := h.service.CreateNote(req.Title, req.Content, orgIDStr, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": note})
}

func (h *NoteHandler) GetAll(c *gin.Context) {
	org, _ := c.Get("orgID")
	orgIDStr := org.(string)
	if orgIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID tidak valid"})
		return
	}

	notes, err := h.service.GetNotes(orgIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notes})
}

func (h *NoteHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus angka"})
		return
	}

	org, _ := c.Get("orgID")
	orgIDStr := org.(string)
	if orgIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID tidak valid"})
		return
	}

	note, err := h.service.GetNoteByID(id, orgIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if note == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan (atau Anda tidak punya akses)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": note})
}
