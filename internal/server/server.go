package server

import (
	"fmt"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/localdrop/internal/config"
	"github.com/yourusername/localdrop/pkg/fileutil"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
	// Ensure upload directory exists
	if err := fileutil.EnsureDir(cfg.UploadDir); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()
	
	// Add recovery middleware
	router.Use(gin.Recovery())

	server := &Server{
		config: cfg,
		router: router,
	}

	// Setup routes
	server.setupRoutes()

	return server, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.printStartupBanner()

	addr := fmt.Sprintf(":%d", s.config.Port)
	if err := s.router.Run(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// printStartupBanner displays server information
func (s *Server) printStartupBanner() {
	localIP := getLocalIP()

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║              LocalShare Server Started                      ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Local:    http://localhost:%d                          ║\n", s.config.Port)
	if localIP != "" {
		fmt.Printf("║  Network:  http://%-15s:%d                      ║\n", localIP, s.config.Port)
	}
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Upload Directory: %-39s ║\n", truncateString(s.config.UploadDir, 39))
	
	if s.config.IsPINProtected() {
		fmt.Println("║  PIN Protection: ENABLED                                ║")
	}
	
	if s.config.IsAdminAuthEnabled() {
		fmt.Println("║  Admin Auth: ENABLED                                    ║")
	}
	
	fmt.Printf("║  Max File Size: %d MB                                  ║\n", s.config.MaxFileSizeMB)
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println("\n Scan this QR code on your phone (or type the Network URL)")
	fmt.Println("Press Ctrl+C to stop the server\n")
}

// getLocalIP returns the local IP address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}