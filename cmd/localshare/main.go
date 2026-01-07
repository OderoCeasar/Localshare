package main

import (
	"fmt"
	"os"

	"github.com/OderoCeasar/localshare/internal/config"
	"github.com/OderoCeasar/localshare/internal/server"
	"github.com/spf13/cobra"
)

func main() {
	var cfg config.Config

	rootCmd := &cobra.Command{
		Use:   "localshare",
		Short: "LocalShare - Simple file sharing on your local network",
		Long: `LocalShare is a cross-platform file-sharing tool that lets you
easily transfer files between devices on your local network via a web browser.

Perfect for quickly transferring files between your laptop and phone, or sharing
files with teammates on the same network without cloud services.`,
		Example: `  # Start with default settings
  LocalShare

  # Start with PIN protection
  LocalShare --pin 1234

  # Start with admin authentication
  LocalShare --admin --admin-pass secret123

  # Custom port and directory
  LocalShare --port 3000 --dir ~/my-shares`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(&cfg)
		},
	}

	// Define flags
	rootCmd.Flags().IntVarP(&cfg.Port, "port", "p", 8080, "Port to run the server on")
	rootCmd.Flags().StringVarP(&cfg.UploadDir, "dir", "d", "./uploads", "Directory to store uploaded files")
	rootCmd.Flags().StringVar(&cfg.PIN, "pin", "", "Optional PIN for file access (4-6 digits)")
	rootCmd.Flags().BoolVar(&cfg.AdminAuth, "admin", false, "Enable admin authentication for uploads")
	rootCmd.Flags().StringVar(&cfg.AdminUser, "admin-user", "admin", "Admin username (when --admin is enabled)")
	rootCmd.Flags().StringVar(&cfg.AdminPass, "admin-pass", "", "Admin password (required when --admin is enabled)")
	rootCmd.Flags().Int64Var(&cfg.MaxFileSizeMB, "max-size", 500, "Maximum file size in MB")

	// Add validation
	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return cfg.Validate()
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runServer(cfg *config.Config) error {
	srv, err := server.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	return srv.Start()
}
