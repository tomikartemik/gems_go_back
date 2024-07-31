package service

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"sync"
	"time"
)

type DropService struct {
	repo repository.Drop
}

type DropResponse struct {
	Photo string `json:"photo_link"`
	Color string `json:"color"`
}

var itemsIds []int

func NewDropService(repo repository.Drop) *DropService {
	itemsIds, _ = repo.GetItemsIds()
	return &DropService{repo: repo}
}

type ClientDrop struct {
	conn *websocket.Conn
}

type InitDropsResponse struct {
	Color string `json:"color"`
	Photo string `json:"photo"`
}

var lastDrops = []DropResponse{}
var clientsMutexDrop = &sync.Mutex{}
var clientsDrop = make(map[*ClientDrop]bool)

func (s *DropService) EditConnsDrop(conn *websocket.Conn) {

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

func (s *DropService) NewDrop(itemId int, dirty bool) {
	item, err := s.repo.NewDrop(itemId, dirty)
	if err == nil {
		if len(lastDrops) >= 7 {
			lastDrops = lastDrops[1:]
		}
		lastDrops = append(lastDrops, DropResponse{Photo: item.PhotoLink, Color: item.Color})
	}
	s.DropWS()
}

func (s *DropService) DirtyMoves() {
	rand.Seed(time.Now().UnixNano())
	for {
		randIndex := rand.Intn(len(itemsIds))
		s.NewDrop(itemsIds[randIndex], true)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func (s *DropService) GetLastDrops() ([]model.Item, error) {
	return s.repo.GetLastDrops()
}
