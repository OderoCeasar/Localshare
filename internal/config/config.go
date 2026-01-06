package config

import (
	"errors"
	"fmt"
	"regexp"
)

type Config struct {
	Port          int
	UploadDir     string
	PIN           string
	AdminAuth     bool
	AdminUser     string
	AdminPass     string
	MaxFileSizeMB int64
}

// Validate checks common CLI-config invariants
func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if c.UploadDir == "" {
		return fmt.Errorf("upload directory must be set")
	}
	if c.PIN != "" {
		// PIN must be 4-6 digits
		matched, _ := regexp.MatchString(`^\d{4,6}$`, c.PIN)
		if !matched {
			return errors.New("pin must be 4-6 digits")
		}
	}
	if c.AdminAuth && c.AdminPass == "" {
		return errors.New("admin authentication enabled but admin password is empty")
	}
	if c.MaxFileSizeMB <= 0 {
		return errors.New("max file size must be greater than 0")
	}
	return nil
}

func (c *Config) IsPINProtected() bool {
	return c.PIN != ""
}

func (c *Config) IsAdminAuthEnabled() bool {
	return c.AdminAuth
}

// MaxFileSize returns the max size in bytes
func (c *Config) MaxFileSize() int64 {
	return c.MaxFileSizeMB * 1024 * 1024
}
