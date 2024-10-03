package handler

import (
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreateFakeBetter(c *gin.Context) {
	var input model.FakeBets
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	h.services.CreateFakeBetter(input)
	c.JSON(http.StatusOK, "OK")
}
