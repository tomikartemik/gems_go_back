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
	Length          float64 `json:"length"`
	Rotate          float64 `json:"rotate"`
}

var respone = Responce{"Crashed", 1.0, 10.0, 0.0, 0.0}
var clients = make(map[*Client]bool)
var clientsMutex = &sync.Mutex{}
var winMultiplier = 0.0
var u = 0.0
var delta = 0.0
var deltaCrash = 0.0

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
	respone.Length = 0.0
	respone.Rotate = 0.0
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
		if respone.Length < 100.0 {
			respone.Length += 0.4
		} else if respone.Rotate < 19.5 {
			respone.Length += 0.0026
			respone.Rotate += 0.0065
		}
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
	delta = respone.Rotate / 300.0
	deltaCrash = 0
	if respone.Length > 100 {
		deltaCrash = (respone.Length - 100) / 300
	}
	for time_before_pending := 300; time_before_pending >= 0; time_before_pending-- {
		respone.Rotate -= delta
		respone.Length -= deltaCrash
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
