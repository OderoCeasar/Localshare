package handlers

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/OderoCeasar/localshare/internal/config"
	"github.com/OderoCeasar/localshare/internal/models"
)

const (
	sessionKeyPIN   = "pin_verified"
	sessionKeyAdmin = "admin_authenticated"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	config *config.Config
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		config: cfg,
	}
}

// VerifyPIN handles PIN verification requests
func (h *AuthHandler) VerifyPIN(c *gin.Context) {
	var req models.PINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(req.PIN), []byte(h.config.PIN)) == 1 {
		session := sessions.Default(c)
		session.Set(sessionKeyPIN, true)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to save session",
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Success: true,
			Message: "PIN verified successfully",
		})
	} else {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid PIN",
		})
	}
}

// AdminLogin handles admin login requests
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	// Use constant-time comparison for both username and password
	userMatch := subtle.ConstantTimeCompare([]byte(req.Username), []byte(h.config.AdminUser)) == 1
	passMatch := subtle.ConstantTimeCompare([]byte(req.Password), []byte(h.config.AdminPass)) == 1

	if userMatch && passMatch {
		session := sessions.Default(c)
		session.Set(sessionKeyAdmin, true)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to save session",
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Success: true,
			Message: "Admin login successfully",
		})
	} else {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid credentials",
		})
	}
}

// AdminLogout handles admin logout requests
func (h *AuthHandler) AdminLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(sessionKeyAdmin)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to save session",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}