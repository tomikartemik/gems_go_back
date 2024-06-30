package handler

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/schema"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary SignUp
// @Tags auth
// @Description Создать нового пользвателя
// @ID signUp
// @Accept json
// @Produce json
// @Param input body schema.InputUser true "User data"
// @Success 200 {object} map[string]interface{} "id"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description Залогиниться
// @ID signIn
// @Accept json
// @Produce json
// @Param input body handler.signInInput true "Credentials"
// @Success 200 {object} map[string]interface{} "token"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.GenerateToken(input.Email, input.Password)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.services.SignIn(input.Email, input.Password)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// @Summary UpdateUser
// @Tags auth
// @Description Обновить инфу о юзере
// @ID updateUser
// @Accept json
// @Produce json
// @Param id query int true "ID юзера"
// @Param input body model.User true "Updates"
// @Success 200 {object} schema.InputUser "UserInfo"
// @Failure 400 {object} error "error"
// @Failure 500 {object} error "error"
// @Router /auth/update [patch]
func (h *Handler) updateUser(c *gin.Context) {
	id := c.Query("id")
	var input schema.InputUser
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	updatedUser, err := h.services.UpdateUser(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) getUserById(c *gin.Context) {
	id := c.Query("id")
	user, err := h.services.GetUserById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) sellAnItem(c *gin.Context) {
	userId := c.Query("user_id")
	userItemIdStr := c.Query("user_item_id")
	var userItemId, err = strconv.Atoi(userItemIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user, err := h.services.SellAnItem(userId, userItemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, user)
}
