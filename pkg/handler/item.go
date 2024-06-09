package handler

import (
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// createItem создает новый элемент.
// @Summary Создает новый элемент
// @Tags items
// @Description Создает новый элемент на основе переданных данных
// @ID createItem
// @Accept json
// @Produce json
// @Param input body model.Item true "Данные для создания элемента"
// @Success 200 {object} map[string]interface{} "Успешно созданный элемент"
// @Failure 400 {object} map[string]interface{} "Ошибка в запросе"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /item/create [post]
func (h *Handler) createItem(c *gin.Context) {
	var input model.Item
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.services.Item.CreateItem(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// getItem получает информацию об элементе.
// @Summary Получает информацию об элементе
// @Tags items
// @Description Получает информацию об элементе по его ID
// @ID getItem
// @Accept json
// @Produce json
// @Param id query int true "ID элемента"
// @Success 200 {object} map[string]interface{} "Информация об элементе"
// @Failure 400 {object} map[string]interface{} "Ошибка в запросе"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /item/get-item [get]
func (h *Handler) getItem(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}

	item, err := h.services.GetItem(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

// getAllItems получает информацию обо всех элементах.
// @Summary Получает информацию обо всех элементах
// @Tags items
// @Description Получает информацию обо всех элементах
// @ID getAllItems
// @Accept json
// @Produce json
// @Success 200 {array} model.ItemWithID "Информация обо всех элементах"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /item/get-all-items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	var items []model.Item
	items, err := h.services.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	var itemsWithId []model.ItemWithID
	for _, item := range items {
		itemsWithId = append(itemsWithId, model.ItemWithID{
			ID:     item.ID,
			Name:   item.Name,
			Rarity: item.Rarity,
			Price:  item.Price,
		})
	}
	c.JSON(http.StatusOK, itemsWithId)
}

// updateItem обновляет информацию об элементе.
// @Summary Обновляет информацию об элементе
// @Tags items
// @Description Обновляет информацию об элементе
// @ID updateItem
// @Accept json
// @Produce json
// @Param input body model.ItemWithID true "Данные для обновления элемента"
// @Success 200 {object} model.ItemWithID "Обновленный элемент"
// @Failure 400 {object} map[string]interface{} "Ошибка в запросе"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /item/update [patch]
func (h *Handler) updateItem(c *gin.Context) {
	var input model.ItemWithID

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	item, err := h.services.Item.UpdateItem(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, item)
}
