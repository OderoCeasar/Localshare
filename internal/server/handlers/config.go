package handlers

import (
	"net/http"

	"github.com/OderoCeasar/localshare/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/OderoCeasar/localshare/internal/config"
	
)


// handles configuration requests
type ConfigHandler struct {
	config *config.Config
}


// NewConfigHandler
func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{
		config: cfg,
	}
}


// returns the server configuration
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, models.ConfigResponse{
		PINProtected: h.config.IsPINProtected(),
		AdminRequired: h.config.IsAdminAuthEnabled(),
		MaxFileSize: h.config.MaxFileSize(),
	})
}