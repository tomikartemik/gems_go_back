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
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, "OK")
}

func (h *Handler) GetReplenishments(c *gin.Context) {
	sortOrder := c.Query("sort_order")
	page := c.DefaultQuery("page", "0")
	pageInt, pageErr := strconv.Atoi(page)
	if pageErr != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, pageErr.Error())
		return
	}
	replenishments, err := h.services.GetReplenishments(sortOrder, pageInt)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, replenishments)
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	var input ChangeStatusRequest

	if role := c.GetString("role"); role != "admin" {
		newErrorResponse(c, http.StatusForbidden, "You are not authorized to access this resource") // 403 Forbidden
		return
	}

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
