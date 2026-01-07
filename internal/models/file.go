package models

import "time"

// FileInfo represents metadata about a file
type FileInfo struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modifiedTime"`
	IsDir        bool      `json:"isDir"`
}

// PINRequest represents a PIN verification request
type PINRequest struct {
	PIN string `json:"pin" binding:"required"`
}

// AdminLoginRequest represents an admin login request
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ConfigResponse represents the server configuration exposed to clients
type ConfigResponse struct {
	PINProtected  bool  `json:"pinProtected"`
	AdminRequired bool  `json:"adminRequired"`
	MaxFileSize   int64 `json:"maxFileSize"`
}

// UploadResponse represents the response after a successful upload
type UploadResponse struct {
	Message  string `json:"message"`
	Filename string `json:"filename"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// FilesListResponse represents a list of files
type FilesListResponse struct {
	Files []FileInfo `json:"files"`
}