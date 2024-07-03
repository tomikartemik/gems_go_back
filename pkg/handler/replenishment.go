package handler

import (
	"fmt"
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) NewReplenishment(c *gin.Context) {
	var input model.Replenishment
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := input.UserID
	amount := input.Amount
	location, err := h.services.Replenishment.NewReplenishment(userID, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, location)
}

func (h *Handler) MSGFromFrekassa(c *gin.Context) {
	var data map[string]interface{}

	// Парсинг JSON из тела запроса
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(data)
	// Для демонстрации выводим полученные данные в консоль
	//c.JSON(http.StatusOK, data)
}
