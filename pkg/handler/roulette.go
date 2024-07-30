package handler

import (
	"gems_go_back/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgraderRoulette = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Разрешить все источники для тестирования
		return true
	},
}

func (h *Handler) handleConnectionsRoulette(c *gin.Context) {
	conn, err := upgraderRoulette.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	h.services.EidtConnsRoulette(conn)
}

func (h *Handler) getAllRouletteRecords(c *gin.Context) {
	var allRecords []model.RouletteRecord
	allRecords, err := h.services.GetAllRouletteRecords()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allRecords)
}

func (h *Handler) initRouletteBetsForNewClient(c *gin.Context) {
	betsAtCurrentGame := h.services.InitRouletteBetsForNewClient()
	c.JSON(http.StatusOK, betsAtCurrentGame)
}
