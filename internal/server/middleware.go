package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/OderoCeasar/localshare/internal/models"
)

const (
	sessionName         = "localshare_session"
	sessionKeyPIN       = "pin_verified"
	sessionKeyAdmin     = "admin_authenticated"
	sessionSecret       = "secret_key"
)

// setupMiddleware configures all middleware for the router
func (s *Server) setupMiddleware() {
	// CORS middleware
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Session middleware
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, 
	})
	s.router.Use(sessions.Sessions(sessionName, store))

	// Custom logger middleware
	s.router.Use(s.loggerMiddleware())
}

// pinMiddleware checks if PIN is verified when PIN protection is enabled
func (s *Server) pinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If PIN protection is disabled, skip verification
		if !s.config.IsPINProtected() {
			c.Next()
			return
		}

		session := sessions.Default(c)
		pinVerified := session.Get(sessionKeyPIN)

		if pinVerified != true {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "PIN verification required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// adminMiddleware checks if admin is authenticated when admin auth is enabled
func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If admin auth is disabled, skip verification
		if !s.config.IsAdminAuthEnabled() {
			c.Next()
			return
		}

		session := sessions.Default(c)
		adminAuth := session.Get(sessionKeyAdmin)

		if adminAuth != true {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Admin authentication required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// loggerMiddleware provides custom logging for requests
func (s *Server) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Build log message
		if raw != "" {
			path = path + "?" + raw
		}

		// Only log non-successful requests and uploads/downloads for monitoring
		if statusCode >= 400 || c.Request.Method == "POST" || c.Request.Method == "DELETE" {
			clientIP := c.ClientIP()
			method := c.Request.Method

			// Color code based on status
			var statusColor string
			switch {
			case statusCode >= 500:
				statusColor = "\033[31m" // Red
			case statusCode >= 400:
				statusColor = "\033[33m" // Yellow
			case statusCode >= 300:
				statusColor = "\033[36m" // Cyan
			default:
				statusColor = "\033[32m" // Green
			}
			resetColor := "\033[0m"

			gin.DefaultWriter.Write([]byte(
				statusColor + "[LocalDrop] " + resetColor +
					method + " " + path +
					" | " + statusColor + string(rune(statusCode)) + resetColor +
					" | " + latency.String() +
					" | " + clientIP + "\n",
			))
		}
	}
}