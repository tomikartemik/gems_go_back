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
	I      int     `json:"i"`
	Ip     string  `json:"ip"`
}

func (h *Handler) NewReplenishment(c *gin.Context) {
	var input ReplInput
	userID := c.GetString("user_id")
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("input: %+v\n", input)
	location, err := h.services.Replenishment.NewReplenishment(userID, input.Amount, input.Promo, input.I, input.Ip)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, location)
}

func (h *Handler) RedirectAccepted(c *gin.Context) {
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

func (h *Handler) MSGFromFreekassa(c *gin.Context) {
	replenishmentIdStr := c.Query("MERCHANT_ORDER_ID")
	replenishmentId, _ := strconv.Atoi(replenishmentIdStr)
	go h.services.CheckReplenishment(replenishmentId)

	c.JSON(http.StatusOK, "OK")
}
