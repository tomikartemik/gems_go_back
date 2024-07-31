package handler

import (
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createWithdraw(c *gin.Context) {
	var input model.Withdraw
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.UserId = c.GetString("user_id")
	err := h.services.CreateWithdraw(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, "OK!")
}

func (h *Handler) getUsersWithdraws(c *gin.Context) {
	userId := c.GetString("user_id")
	h.services.Withdraw.GetUsersWithdraws(userId)
}

func (h *Handler) getPositionPrices(c *gin.Context) {
	c.JSON(http.StatusOK, h.services.GetPositionPrices())
}
