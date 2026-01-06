package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/OderoCeasar/localshare/internal/server/handlers"
)

// setupRoutes configures all routes for the application
func (s *Server) setupRoutes() {
	// Setup middleware first
	s.setupMiddleware()

	// Create handlers
	authHandler := handlers.NewAuthHandler(s.config)
	fileHandler := handlers.NewFileHandler(s.config)
	configHandler := handlers.NewConfigHandler(s.config)

	// Serve static frontend
	s.router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(getIndexHTML()))
	})

	// API routes
	api := s.router.Group("/api")
	{
		// Public endpoints
		api.GET("/config", configHandler.GetConfig)
		api.POST("/verify-pin", authHandler.VerifyPIN)
		api.POST("/admin/login", authHandler.AdminLogin)
		api.POST("/admin/logout", authHandler.AdminLogout)

		// Protected file endpoints (require PIN if enabled)
		files := api.Group("/files")
		files.Use(s.pinMiddleware())
		{
			files.GET("", fileHandler.ListFiles)
			files.GET("/download/:filename", fileHandler.DownloadFile)
			
			// These also require admin auth if enabled
			files.POST("/upload", s.adminMiddleware(), fileHandler.UploadFile)
			files.DELETE("/:filename", s.adminMiddleware(), fileHandler.DeleteFile)
		}
	}

	// Health check endpoint
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "localshare",
		})
	})
}

