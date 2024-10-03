package handler

import (
	"gems_go_back/pkg/model"
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

func (h *Handler) signUpAdmin(c *gin.Context) {
	var input model.Admin

	if role := c.GetString("role"); role != "admin" {
		newErrorResponse(c, http.StatusForbidden, "You are not authorized to access this resource") // 403 Forbidden
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Admin.CreateAdmin(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func (h *Handler) signInAdmin(c *gin.Context) {
	var input model.Admin

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Admin.SignInAdmin(input.Username, input.Password)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
