package handler

import (
	"fmt"
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) NewReplenishment(c *gin.Context) {
	var input model.Replenishment
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userID := input.UserID
	amount := input.Amount
	location, err := h.services.Replenishment.NewReplenishment(userID, amount)
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
