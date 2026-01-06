package fileutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/OderoCeasar/localshare/internal/models"
)

// ErrInvalidPath is returned when a path contains invalid characters
var ErrInvalidPath = errors.New("invalid file path")

// SanitizeFilename removes 
func SanitizeFilename(filename string) (string, error) {

	base := filepath.Base(filename)
	
	// Check for path traversal attempts
	if strings.Contains(base, "..") || strings.Contains(base, "/") || strings.Contains(base, "\\") {
		return "", ErrInvalidPath
	}
	
	// Check for empty filename
	if base == "" || base == "." {
		return "", ErrInvalidPath
	}
	
	return base, nil
}

// ListFiles returns a list of files in the specified directory
func ListFiles(dirPath string) ([]models.FileInfo, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	files := make([]models.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			// Skip files we can't read
			continue
		}

		files = append(files, models.FileInfo{
			Name:         entry.Name(),
			Size:         info.Size(),
			ModifiedTime: info.ModTime(),
			IsDir:        entry.IsDir(),
		})
	}

	return files, nil
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// DeleteFile removes a file from the filesystem
func DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %w", err)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFilePath safely joins the directory and filename
func GetFilePath(dir, filename string) (string, error) {
	sanitized, err := SanitizeFilename(filename)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, sanitized), nil
}