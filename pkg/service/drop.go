package service

import (
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

func (s *DropService) GetLastDrops() ([]model.Item, error) {
	return s.repo.GetLastDrops()
}

type ClientDrop struct {
	conn *websocket.Conn
	send chan []model.Item
}

var lastDrops = []model.Item{}
var clientsMutexDrop = &sync.Mutex{}
var clientsDrop = make(map[*ClientDrop]bool)

func (s *DropService) EidtConnsDrop(conn *websocket.Conn) {
	client := &ClientDrop{conn: conn, send: make(chan []model.Item)}

	clientsMutexDrop.Lock()
	clientsDrop[client] = true
	clientsMutexDrop.Unlock()

	go s.handleClient(client)
}

func (s *DropService) handleClient(client *ClientDrop) {
	defer func() {
		clientsMutexDrop.Lock()
		delete(clientsDrop, client)
		clientsMutexDrop.Unlock()
		client.conn.Close()
	}()

	for msg := range client.send {
		err := client.conn.WriteJSON(msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func (s *DropService) DropsWS() {
	clientsMutexDrop.Lock()
	defer clientsMutexDrop.Unlock()

	for client := range clientsDrop {
		select {
		case client.send <- lastDrops:
		default:
			close(client.send)
			delete(clientsDrop, client)
		}
	}
}

func (s *DropService) NewDrop(itemId int) {
	item, err := s.repo.NewDrop(itemId)
	if err == nil {
		if len(lastDrops) >= 7 {
			lastDrops = lastDrops[1:]
		}
		lastDrops = append(lastDrops, item)
	}
	s.DropsWS()
}
