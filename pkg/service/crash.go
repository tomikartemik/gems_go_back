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

type ClientCrash struct {
	conn *websocket.Conn
}

type ResponseCrash struct {
	Status          string  `json:"status"`
	Multiplier      float64 `json:"multiplier"`
	TimeBeforeStart float64 `json:"time_before_start"`
	Length          float64 `json:"length"`
	Rotate          float64 `json:"rotate"`
}

var responseCrash = ResponseCrash{"Crashed", 1.0, 10.0, 0.0, 0.0}
var clientsCrash = make(map[*ClientCrash]bool)
var clientsMutexCrash = &sync.Mutex{}
var winMultiplier = 0.0
var u = 0.0
var delta = 0.0
var deltaCrash = 0.0

var stepen = math.Pow(2.0, 52.0)

func (s *CrashService) EditConnsCrash(conn *websocket.Conn) {

	defer conn.Close()

	client := &ClientCrash{conn: conn}
	clientsMutexCrash.Lock()
	clientsCrash[client] = true
	clientsMutexCrash.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
	}

	clientsMutexCrash.Lock()
	delete(clientsCrash, client)
	clientsMutexCrash.Unlock()
}

func (s *CrashService) BroadcastTimeCrash() {
	startPreparingCrash()
}

func startPreparingCrash() {
	responseCrash.Length = 0.0
	responseCrash.Rotate = 0.0
	responseCrash.Status = "Pending"
	u = rand.Float64() * (stepen)
	winMultiplier = math.Round((100*stepen-u)/(stepen-u)) / 100.0
	preparingCrash()
}
func preparingCrash() {
	for time_before_start := 1000.0; time_before_start >= 0; time_before_start-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutexCrash.Lock()
		responseCrash.TimeBeforeStart = time_before_start / 100.0
		for client := range clientsCrash {
			err := client.conn.WriteJSON(responseCrash)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsCrash, client)
			}
		}
		clientsMutexCrash.Unlock()
	}
	startGameCrash()
}

func startGameCrash() {
	responseCrash.Status = "Running"
	responseCrash.Multiplier = 1.0
	gameCrash()
}

func gameCrash() {
	for responseCrash.Multiplier < winMultiplier {
		time.Sleep(10 * time.Millisecond)
		//responseCrash.Multiplier = responseCrash.Multiplier * 1.0004
		responseCrash.Multiplier = math.Round(responseCrash.Multiplier*10003) / 10000
		if responseCrash.Length < 100.0 {
			responseCrash.Length += 0.4
		} else if responseCrash.Rotate < 19.5 {
			responseCrash.Length += 0.0026
			responseCrash.Rotate += 0.0065
		}
		clientsMutexCrash.Lock()
		for client := range clientsCrash {
			err := client.conn.WriteJSON(responseCrash)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsCrash, client)
			}
		}
		clientsMutexCrash.Unlock()
	}
	endCrash()
}

func endCrash() {
	responseCrash.Status = "Crashed"
	delta = responseCrash.Rotate / 300.0
	deltaCrash = 0
	if responseCrash.Length > 100 {
		deltaCrash = (responseCrash.Length - 100) / 300
	}
	for time_before_pending := 300; time_before_pending >= 0; time_before_pending-- {
		responseCrash.Rotate -= delta
		responseCrash.Length -= deltaCrash
		time.Sleep(10 * time.Millisecond)
		clientsMutexCrash.Lock()
		for client := range clientsCrash {
			err := client.conn.WriteJSON(responseCrash)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsCrash, client)
			}
		}
		clientsMutexCrash.Unlock()
	}
	startPreparingCrash()
}
