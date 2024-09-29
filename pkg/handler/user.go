package handler

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/schema"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	signInResponse, err := h.services.SignIn(input.Email, input.Password)

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, signInResponse)
	h.services.Online.SetOnline()
}

func (h *Handler) updateUser(c *gin.Context) {
	id := c.GetString("user_id")
	var input schema.InputUser
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	updatedUser, err := h.services.UpdateUser(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) getUserById(c *gin.Context) {
	id := c.GetString("user_id")
	user, err := h.services.GetUserById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) sellItem(c *gin.Context) {
	userId := c.GetString("user_id")
	userItemIdStr := c.Query("user_item_id")
	var userItemId, err = strconv.Atoi(userItemIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := h.services.SellItem(userId, userItemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) sellAllItems(c *gin.Context) {
	userId := c.GetString("user_id")
	err := h.services.SellAllItems(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "OK!")
}

func (h *Handler) changeAvatar(c *gin.Context) {
	userId := c.GetString("user_id")
	newPhoto, err := h.services.ChangeAvatar(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "не хватает денежек")
		return
	}
	c.JSON(http.StatusOK, newPhoto)
	h.services.Online.SetOnline()
}
