package models

import "time"

type FileInfo struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	IsDir    bool      `json:"is_dir"`
	Modified time.Time `json:"modified"`
}

type PINRequest struct {
	PIN string `json:"pin"`
}

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UploadResponse struct {
	Message  string `json:"message"`
	Filename string `json:"filename"`
}

type FilesListResponse struct {
	Files []FileInfo `json:"files"`
}

type ConfigResponse struct {
	PINProtected  bool  `json:"pin_protected"`
	AdminRequired bool  `json:"admin_required"`
	MaxFileSize   int64 `json:"max_file_size"`
}
