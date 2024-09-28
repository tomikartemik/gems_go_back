package handler

import (
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) UpdateCard(c *gin.Context) {
	var input model.Card
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err := h.services.UpdateCard(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, "OK")
}

func (h *Handler) GetCard(c *gin.Context) {
	card, err := h.services.GetCard()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, card)
}
