package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ChangeStatusRequest struct {
	ReplenishmentId int    `json:"replenishment_id"`
	Status          string `json:"status"`
}

func (h *Handler) CreateOwnReplenishment(c *gin.Context) {
	userId := c.GetString("user_id")
	amountStr := c.Query("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.CreateReplenishment(amount, userId, file)
}

func (h *Handler) GetReplenishments(c *gin.Context) {
	replenishments, err := h.services.GetReplenishments()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, replenishments)
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	var input ChangeStatusRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err := h.services.ChangeStatus(input.ReplenishmentId, input.Status)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "OK")
}
