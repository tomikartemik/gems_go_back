package handler

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/schema"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OpenedCaseResponse struct {
	WinedItem  model.ItemWithID `json:"wined_item"`
	UserItemId int              `json:"user_item_id"`
}

func (h *Handler) createCase(c *gin.Context) {
	var input model.Case
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	createdCase, err := h.services.CreateCase(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, createdCase)
}

func (h *Handler) getCase(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
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
		return
	}
	h.services.Online.SetOnline()
	c.JSON(http.StatusOK, cases)
}

func (h *Handler) updateCase(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input model.Case
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	updatedCaseInfo, err := h.services.UpdateCase(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, updatedCaseInfo)
}

func (h *Handler) deleteCase(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.DeleteCase(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) openCase(c *gin.Context) {
	userId := c.GetString("user_id")
	caseIdStr := c.Query("case_id")
	caseId, err := strconv.Atoi(caseIdStr)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	chosenItem, userItemId, err := h.services.OpenCase(userId, caseId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else if userItemId == -1 {
		c.JSON(http.StatusBadRequest, "иди у мамы денег проси")
		return
	} else {
		openedCaseResponse := OpenedCaseResponse{
			WinedItem:  chosenItem,
			UserItemId: userItemId,
		}
		c.JSON(http.StatusOK, openedCaseResponse)
	}
	h.services.Online.SetOnline()
}

func (h *Handler) getAllCaseRecords(c *gin.Context) {
	var allCaseRecords []schema.CaseInfo
	allCaseRecords, err := h.services.GetAllCaseRecords()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allCaseRecords)
}

func (h *Handler) getLastDrops(c *gin.Context) {
	lastDrops, err := h.services.GetLastDrops()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		c.JSON(http.StatusOK, lastDrops)
	}
}
