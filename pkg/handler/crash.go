package handler

import (
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgraderCrash = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Разрешить все источники для тестирования
		return true
	},
}

func (h *Handler) handleConnectionsCrash(c *gin.Context) {
	conn, err := upgraderCrash.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	h.services.EditConnsCrash(conn)
}

func (h *Handler) getAllCrashRecords(c *gin.Context) {
	var allRecords []model.CrashRecord
	allRecords, err := h.services.GetAllRecords()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, allRecords)
}

func (h *Handler) initCrashBetsForNewClient(c *gin.Context) {
	betsAtCurrentGame := h.services.InitCrashBetsForNewClient()
	c.JSON(http.StatusOK, betsAtCurrentGame)
}
