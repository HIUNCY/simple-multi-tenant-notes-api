package handler

import (
	"net/http"

	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/config"
	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type loginRequest struct {
	UserID string `json:"user_id" binding:"required"`
	OrgID  string `json:"org_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(req.UserID, req.OrgID, req.Role, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login berhasil! Gunakan token ini di Header 'Authorization: Bearer <token>'",
	})
}
