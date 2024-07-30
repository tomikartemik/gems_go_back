package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AdminChangeStatus struct {
	Status bool `json:"status"`
}

func (h *Handler) adminChangeStatus(c *gin.Context) {
	var input AdminChangeStatus
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.services.ChangeStatusOfStartCrash(input.Status)
	h.services.ChangeStatusOfStartRoulette(input.Status)
	h.services.SetOnline()
	c.JSON(http.StatusOK, "Changed status to "+strconv.FormatBool(input.Status))
}
