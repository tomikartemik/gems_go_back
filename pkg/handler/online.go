package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getOnline(c *gin.Context) {
	usersOnline := h.services.Online.GetOnline()
	c.JSON(http.StatusOK, usersOnline)
}
