package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReplInput struct {
	Amount float64 `json:"amount"`
	Promo  string  `json:"promo"`
}

func (h *Handler) NewReplenishment(c *gin.Context) {
	var input ReplInput
	userID := c.GetString("user_id")
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	location, err := h.services.Replenishment.NewReplenishment(userID, input.Amount, input.Promo)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, location)
}

func (h *Handler) RedirectAccepted(c *gin.Context) {
	replenishmentIdStr := c.Query("MERCHANT_ORDER_ID")
	replenishmentId, _ := strconv.Atoi(replenishmentIdStr)
	go h.services.AcceptReplenishment(replenishmentId)
	c.Redirect(http.StatusFound, "https://dododrop.ru")
}

func (h *Handler) RedirectDenied(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	for key, values := range queryParams {
		for _, value := range values {
			c.String(http.StatusOK, "Key: %s, Value: %s\n", key, value)
		}
	}
	c.Redirect(http.StatusFound, "https://dododrop.ru")
}

func (h *Handler) MSGFromFrekassa(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(data)
	c.Redirect(http.StatusFound, "https://dododrop.ru")
}
