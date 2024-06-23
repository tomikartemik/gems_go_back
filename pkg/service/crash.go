package service

import (
	"gems_go_back/pkg/repository"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

type CrashService struct {
	repo repository.Crash
}

func NewCrashService(repo repository.Crash) *CrashService {
	return &CrashService{repo: repo}
}

type Client struct {
	conn *websocket.Conn
}

type Responce struct {
	Status          string  `json:"status"`
	Multiplier      float64 `json:"multiplier"`
	TimeBeforeStart float64 `json:"time_before_start"`
}

var respone = Responce{"Crashed", 1.0, 10.0}
var clients = make(map[*Client]bool)
var clientsMutex = &sync.Mutex{}
var winMultiplier = 0.0
var u = 0.0

var stepen = math.Pow(2.0, 52.0)

func (s *CrashService) EditConns(conn *websocket.Conn) {

	defer conn.Close()

	client := &Client{conn: conn}
	clientsMutex.Lock()
	clients[client] = true
	clientsMutex.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
	}

	clientsMutex.Lock()
	delete(clients, client)
	clientsMutex.Unlock()
}

func (s *CrashService) BroadcastTime() {
	startPreparing()
}

func startPreparing() {
	respone.TimeBeforeStart = 1.0
	respone.Status = "Pending"
	u = rand.Float64() * (stepen)
	winMultiplier = math.Round((100*stepen-u)/(stepen-u)) / 100.0
	preparing()
}
func preparing() {
	for time_before_start := 1000.0; time_before_start >= 0; time_before_start-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutex.Lock()
		respone.TimeBeforeStart = time_before_start / 100.0
		for client := range clients {
			err := client.conn.WriteJSON(respone)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
	startGame()
}

func startGame() {
	respone.Status = "Running"
	respone.Multiplier = 1.0
	game()
}

func game() {
	for respone.Multiplier < winMultiplier {
		time.Sleep(10 * time.Millisecond)
		//respone.Multiplier = respone.Multiplier * 1.0004
		respone.Multiplier = math.Round(respone.Multiplier*10003) / 10000

		clientsMutex.Lock()
		for client := range clients {
			err := client.conn.WriteJSON(respone)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
	end()
}

func end() {
	respone.Status = "Crashed"
	for time_before_pending := 300; time_before_pending >= 0; time_before_pending-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutex.Lock()
		for client := range clients {
			err := client.conn.WriteJSON(respone)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
	startPreparing()
}
