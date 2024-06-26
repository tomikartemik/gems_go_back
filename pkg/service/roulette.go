package service

import (
	"gems_go_back/pkg/repository"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"sync"
	"time"
)

type RouletteService struct {
	repo repository.Roulette
}

func NewRouletteService(repo repository.Roulette) *RouletteService {
	return &RouletteService{repo: repo}
}

type ClientRoulette struct {
	conn *websocket.Conn
}

type ResponseRoulette struct {
	Status          string  `json:"status"`
	Cell            int     `json:"cell"`
	TimeBeforeStart float64 `json:"time_before_start"`
}

type Cell struct {
	Value  int
	Weight int
}

var clientsRoulette = make(map[*ClientRoulette]bool)
var clientsMutexRoulette = &sync.Mutex{}
var responeRoulette = ResponseRoulette{"Pending", 0, 0.0}
var cells = []Cell{
	{2, 50},
	{3, 33},
	{5, 20},
	{10, 10},
	{100, 1},
}
var totalWeight = 114
var winCell = 0

func (s *RouletteService) EidtConnsRoulette(conn *websocket.Conn) {

	defer conn.Close()

	client := &ClientRoulette{conn: conn}
	clientsMutexRoulette.Lock()
	clientsRoulette[client] = true
	clientsMutexRoulette.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
	}

	clientsMutexRoulette.Lock()
	delete(clientsRoulette, client)
	clientsMutexRoulette.Unlock()
}
func (s *RouletteService) BroadcastTimeRoulette() {
	startPreparingRoulette()
}

func startPreparingRoulette() {
	responeRoulette.Cell = 0
	responeRoulette.Status = "Pending"
	responeRoulette.TimeBeforeStart = 10.0
	preparingRoulette()
}

func preparingRoulette() {
	for time_before_start := 1000.0; time_before_start >= 0; time_before_start-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutexRoulette.Lock()
		responeRoulette.TimeBeforeStart = time_before_start / 100.0
		for client := range clientsRoulette {
			err := client.conn.WriteJSON(responeRoulette)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsRoulette, client)
			}
		}
		clientsMutexRoulette.Unlock()
	}
	startGameRoulette()
}

func startGameRoulette() {
	responeRoulette.Status = "Playing"
	randomNumber := rand.Intn(totalWeight)

	for _, choosenCell := range cells {
		if randomNumber < choosenCell.Weight {
			winCell = choosenCell.Value
			break
		}
		randomNumber -= choosenCell.Weight
	}
	responeRoulette.Cell = winCell
	gameRoulette()
}

func gameRoulette() {
	for time_before_end := 1000.0; time_before_end >= 0; time_before_end-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutexRoulette.Lock()
		for client := range clientsRoulette {
			err := client.conn.WriteJSON(responeRoulette)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsRoulette, client)
			}
		}
		clientsMutexRoulette.Unlock()
	}
	endRoulette()
}

func endRoulette() {
	responeRoulette.Status = "End"
	for time_before_pending := 300; time_before_pending >= 0; time_before_pending-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutexRoulette.Lock()
		for client := range clientsRoulette {
			err := client.conn.WriteJSON(responeRoulette)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsRoulette, client)
			}
		}
		clientsMutexRoulette.Unlock()
	}
	startPreparingRoulette()
}
