package handler

import (
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
