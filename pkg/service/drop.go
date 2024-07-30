package service

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type DropService struct {
	repo repository.Drop
}

func NewDropService(repo repository.Drop) *DropService {
	return &DropService{repo: repo}
}

type ClientDrop struct {
	conn *websocket.Conn
}

var lastDrops = []model.Item{}
var clientsMutexDrop = &sync.Mutex{}
var clientsDrop = make(map[*ClientDrop]bool)

func (s *DropService) EidtConnsDrop(conn *websocket.Conn) {

	defer conn.Close()

	client := &ClientDrop{conn: conn}
	clientsMutexDrop.Lock()
	clientsDrop[client] = true
	clientsMutexDrop.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
	}

	clientsMutexDrop.Lock()
	delete(clientsDrop, client)
	clientsMutexDrop.Unlock()
}

func (s *DropService) DropWS() {
	clientsMutexDrop.Lock()
	for client := range clientsDrop {
		err := client.conn.WriteJSON(lastDrops)
		if err != nil {
			log.Println("Write error:", err)
			client.conn.Close()
			delete(clientsDrop, client)
		}
	}
	clientsMutexDrop.Unlock()

}

func (s *DropService) NewDrop(itemId int) {
	item, err := s.repo.NewDrop(itemId)
	if err == nil {
		if len(lastDrops) >= 7 {
			lastDrops = lastDrops[1:]
		}
		lastDrops = append(lastDrops, item)
	}
	s.DropWS()
}

func (s *DropService) GetLastDrops() ([]model.Item, error) {
	return s.repo.GetLastDrops()
}
