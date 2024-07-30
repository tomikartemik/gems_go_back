package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func (h *Handler) getLastDrops(c *gin.Context) {
	drops, err := h.services.GetLastDrops()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, drops)
}

var upgraderDrop = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Разрешить все источники для тестирования
		return true
	},
}

func (h *Handler) handleConnectionsDrop(c *gin.Context) {
	conn, err := upgraderDrop.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	h.services.EidtConnsDrop(conn)
}
