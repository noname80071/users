package files

import (
	"go-users/internal/domain/ports"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ports.FilesServicePort
}

func NewHandler(service ports.FilesServicePort) *Handler {
	return &Handler{service: service}
}

func (h *Handler) UploadSkin(c *gin.Context) {
	userID := c.Param("id")

	file, err := c.FormFile("skin")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get skin file"})
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer fileReader.Close()

	url, err := h.service.UploadSkin(c.Request.Context(), userID, fileReader, file.Filename, file.Size)

	// Обработка ошибок!!!!
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "skin uploaded successfully",
		"user_id": userID,
		"url":     url,
	})
}

func (h *Handler) UploadCloak(c *gin.Context) {
	userID := c.Param("id")

	file, err := c.FormFile("cloak")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get cloak file"})
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer fileReader.Close()

	url, err := h.service.UploadCloak(c.Request.Context(), userID, fileReader, file.Filename, file.Size)

	// Обработка ошибок!!!!
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cloak uploaded successfully",
		"user_id": userID,
		"url":     url,
	})
}
