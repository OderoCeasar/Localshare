package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/OderoCeasar/localshare/internal/config"
	"github.com/OderoCeasar/localshare/internal/models"
	"github.com/OderoCeasar/localshare/pkg/fileutil"
)

// FileHandler handles file-related requests
type FileHandler struct {
	config *config.Config
}

// NewFileHandler creates a new file handler
func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{
		config: cfg,
	}
}

// ListFiles returns a list of all uploaded files
func (h *FileHandler) ListFiles(c *gin.Context) {
	files, err := fileutil.ListFiles(h.config.UploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to list files",
		})
		return
	}

	c.JSON(http.StatusOK, models.FilesListResponse{
		Files: files,
	})
}

// DownloadFile sends a file to the client
func (h *FileHandler) DownloadFile(c *gin.Context) {
	filename := c.Param("filename")

	// Get safe file path
	filePath, err := fileutil.GetFilePath(h.config.UploadDir, filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid filename",
		})
		return
	}

	// Check if file exists
	if !fileutil.FileExists(filePath) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "File not found",
		})
		return
	}

	// Send file
	c.File(filePath)
}

// UploadFile handles file upload requests
func (h *FileHandler) UploadFile(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file provided",
		})
		return
	}

	// Check file size
	if file.Size > h.config.MaxFileSize() {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: fmt.Sprintf(
				"File size exceeds maximum of %d MB",
				h.config.MaxFileSizeMB,
			),
		})
		return
	}

	// Sanitize filename
	safeFilename, err := fileutil.SanitizeFilename(file.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid filename",
		})
		return
	}

	// Build destination path
	dst := filepath.Join(h.config.UploadDir, safeFilename)

	// Save file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to save file",
		})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		Message:  "File uploaded successfully",
		Filename: safeFilename,
	})
}

// DeleteFile removes a file from the server
func (h *FileHandler) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")

	// Get safe file path
	filePath, err := fileutil.GetFilePath(h.config.UploadDir, filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid filename",
		})
		return
	}

	// Delete file
	if err := fileutil.DeleteFile(filePath); err != nil {
		if fileutil.FileExists(filePath) {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to delete file",
			})
		} else {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "File not found",
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "File deleted successfully",
	})
}