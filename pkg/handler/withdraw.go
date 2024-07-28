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
	err := h.services.CreateWithdraw(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, "OK!")
}

func (h *Handler) getUsersWithdraws(c *gin.Context) {
	userId := c.Query("user_id")
	h.services.Withdraw.GetUsersWithdraws(userId)
}
