package handler

import (
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// createCase создает новый кейс.
// @Summary Создает новый кейс
// @Tags cases
// @Description Создает новый кейс на основе переданных данных
// @ID createCase
// @Accept json
// @Produce json
// @Param input body schema.CaseInput true "Данные для создания кейса"
// @Success 200 {object} map[string]interface{} "Успешно созданный кейс"
// @Failure 400 {object} map[string]interface{} "Ошибка в запросе"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /case/create [post]
func (h *Handler) createCase(c *gin.Context) {
	var input model.Case
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCase, err := h.services.CreateCase(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, createdCase)
}

// getCase получает информацию о кейсе.
// @Summary Получает информацию о кейсе
// @Tags cases
// @Description Получает информацию о кейсе по его ID
// @ID getCase
// @Accept json
// @Produce json
// @Param id query int true "ID кейса"
// @Success 200 {object} map[string]interface{} "Информация о кейсе"
// @Failure 400 {object} map[string]interface{} "Ошибка в запросе"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /case/get-case [get]
func (h *Handler) getCase(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	caseInfo, err := h.services.GetCase(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, caseInfo)
}

func (h *Handler) getAllCases(c *gin.Context) {
	cases, err := h.services.GetAllCases()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, cases)
}

func (h *Handler) updateCase(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var input model.Case
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	updatedCaseInfo, err := h.services.UpdateCase(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, updatedCaseInfo)
}

func (h *Handler) deleteCase(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = h.services.DeleteCase(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) openCase(c *gin.Context) {
	caseIdStr := c.Query("case_id")
	caseId, err := strconv.Atoi(caseIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	userId := c.Query("user_id")
	chosenItem, err := h.services.OpenCase(caseId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	_, err = h.services.AddItemToInventory(userId, chosenItem.ID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, chosenItem)
}
