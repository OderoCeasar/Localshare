package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/OderoCeasar/localshare/internal/config"
	"github.com/OderoCeasar/localshare/internal/models"
)

// ConfigHandler handles configuration-related requests
type ConfigHandler struct {
	config *config.Config
}

// NewConfigHandler creates a new config handler
func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{
		config: cfg,
	}
}

// GetConfig returns the server configuration
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, models.ConfigResponse{
		PINProtected:  h.config.IsPINProtected(),
		AdminRequired: h.config.IsAdminAuthEnabled(),
		MaxFileSize:   h.config.MaxFileSize(),
	})
}