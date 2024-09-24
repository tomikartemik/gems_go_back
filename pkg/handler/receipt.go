package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) UploadReceipt(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить файл"})
		return
	}

	fileURL, err := h.services.UploadReceipt(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить файл"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Файл успешно загружен",
		"file_url": fileURL,
	})
}
