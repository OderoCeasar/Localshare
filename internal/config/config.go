package config

import (
	"errors"
	"fmt"
	"regexp"
)

// Config holds all application configuration
type Config struct {
	Port          int
	UploadDir     string
	PIN           string
	AdminAuth     bool
	AdminUser     string
	AdminPass     string
	MaxFileSizeMB int64
}

// MaxFileSize returns the maximum file size in bytes
func (c *Config) MaxFileSize() int64 {
	return c.MaxFileSizeMB * 1024 * 1024
}

// IsPINProtected returns whether PIN protection is enabled
func (c *Config) IsPINProtected() bool {
	return c.PIN != ""
}

// IsAdminAuthEnabled returns whether admin authentication is required
func (c *Config) IsAdminAuthEnabled() bool {
	return c.AdminAuth
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate port
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}

	// Validate PIN format if provided
	if c.PIN != "" {
		if !isValidPIN(c.PIN) {
			return errors.New("PIN must be 4-6 digits")
		}
	}

	// Validate admin configuration
	if c.AdminAuth {
		if c.AdminPass == "" {
			return errors.New("admin password is required when admin authentication is enabled (use --admin-pass)")
		}
		if len(c.AdminPass) < 6 {
			return errors.New("admin password must be at least 6 characters")
		}
	}

	// Validate max file size
	if c.MaxFileSizeMB < 1 {
		return errors.New("max file size must be at least 1 MB")
	}
	if c.MaxFileSizeMB > 10000 {
		return errors.New("max file size cannot exceed 10000 MB (10 GB)")
	}

	return nil
}

// isValidPIN checks if the PIN is 4-6 digits
func isValidPIN(pin string) bool {
	matched, _ := regexp.MatchString(`^\d{4,6}$`, pin)
	return matched
}