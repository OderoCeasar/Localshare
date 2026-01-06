package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/OderoCeasar/localshare/internal/models"
	"github.com/OderoCeasar/localshare/pkg/fileutil"
	"github.com/gin-gonic/gin"
	"github.com/OderoCeasar/localshare/internal/config"
)

// handles file requests
type FileHandler struct {
	config *config.Config
}


// create new file handler
func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{
		config: cfg,
	}
}


// listfiles(returns a list of uploaded files)
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


// downloadfile sends a file to the client
func (h *FileHandler) DownloadFile(c *gin.Context) {
	filename := c.Param("filename")

	// get safe file path
	filepath, err := fileutil.GetFilePath(h.config.UploadDir, filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid filename",
		})
		return
	}


	// checks if file exists
	if !fileutil.FileExists(filepath) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "File not found",
		})
		return
	}

	// send file
	c.File(filepath)
}



// Uploadfile handles file upload requests
func (h *FileHandler) UploadFile(c *gin.Context) {
	// get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file provided",
		})
		return
	}

	// checks file size
	if file.Size > h.config.MaxFileSize() {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: fmt.Sprintf(
				"File size exceeds maximum of %d MB",
				h.config.MaxFileSizeMB,
			),
		})
		return
	}

	safeFilename, err := fileutil.SanitizeFilename(file.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid filename",
		})
		return
	}


	dst := filepath.Join(h.config.UploadDir, safeFilename)

	// save file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to save file",
		})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		Message: "File uploaded successfully",
		Filename: safeFilename,
	})
}


// DeleteFile removes a file form the server
func (h *FileHandler) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")

	filepath, err := fileutil.GetFilePath(h.config.UploadDir, filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid filename",
		})
		return
	}

	// delete file
	if err := fileutil.DeleteFile(filepath); err != nil {
		if fileutil.FileExists(filepath) {
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