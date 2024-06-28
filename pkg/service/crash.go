package service

import (
	"encoding/json"
	"fmt"
	"gems_go_back/pkg/model"
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

type BetMessageCrash struct {
	GameId   int     `json:"game_id"`
	PlayerID string  `json:"player_id"`
	Amount   float64 `json:"amount"`
}

type CashoutMessageCrash struct {
	GameId     int     `json:"game_id"`
	PlayerID   string  `json:"player_id"`
	Multiplier float64 `json:"multiplier"`
}

type ResponseCrash struct {
	GameID          int     `json:"game_id"`
	Status          string  `json:"status"`
	Multiplier      float64 `json:"multiplier"`
	TimeBeforeStart float64 `json:"time_before_start"`
	Length          float64 `json:"length"`
	Rotate          float64 `json:"rotate"`
}

var responseCrash = ResponseCrash{0, "Crashed", 1.0, 10.0, 0.0, 0.0}
var clientsCrash = make(map[*ClientCrash]bool)
var clientsMutexCrash = &sync.Mutex{}
var winMultiplier = 0.0
var u = 0.0
var delta = 0.0
var deltaCrash = 0.0
var lastGameID int

var stepen = math.Pow(2.0, 52.0)
var acceptingBetsCrash = true

var acceptingCashoutsCrash = false

func (s *CrashService) EditConnsCrash(conn *websocket.Conn) {

	defer conn.Close()

	client := &ClientCrash{conn: conn}
	clientsMutexCrash.Lock()
	clientsCrash[client] = true
	clientsMutexCrash.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		if acceptingBetsCrash {
			var bet BetMessageCrash
			if err = json.Unmarshal(message, &bet); err != nil {
				log.Println("Invalid bet format:", err)
				continue
			}
			newBet := model.BetCrash{
				GameId: bet.GameId,
				UserID: bet.PlayerID,
				Amount: bet.Amount,
			}
			errorStr := s.repo.NewBetCrash(newBet)
			fmt.Println(errorStr)
		} else if acceptingCashoutsCrash {
			var cashout CashoutMessageCrash
			if err = json.Unmarshal(message, &cashout); err != nil {
				log.Println("Invalid bet format:", err)
				continue
			}
			errorStr := s.repo.NewCashoutCrash(cashout.GameId, cashout.PlayerID, cashout.Multiplier)
			fmt.Println(errorStr)
		}
	}

	clientsMutexCrash.Lock()
	delete(clientsCrash, client)
	clientsMutexCrash.Unlock()
}

func (s *CrashService) BroadcastTimeCrash() {
	s.StartPreparingCrash()
}

func (s *CrashService) StartPreparingCrash() {
	acceptingBetsCrash = true
	responseCrash.Length = 0.0
	responseCrash.Rotate = 0.0
	responseCrash.Status = "Pending"
	u = rand.Float64() * (stepen)
	winMultiplier = math.Round((100*stepen-u)/(stepen-u)) / 100.0
	lastGame, err := s.repo.GetLastRecord()
	if err != nil {
		log.Fatal(err)
	}
	lastGameID = lastGame.ID + 1
	responseCrash.GameID = lastGameID
	s.PreparingCrash()
}

func (s *CrashService) PreparingCrash() {
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
	s.StartGameCrash()
}

func (s *CrashService) StartGameCrash() {
	acceptingBetsCrash = false
	acceptingCashoutsCrash = true
	responseCrash.Status = "Running"
	responseCrash.Multiplier = 1.0
	s.GameCrash()
}

func (s *CrashService) GameCrash() {
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
	go s.repo.NewRecord(winMultiplier)
	s.EndCrash()
}

func (s *CrashService) EndCrash() {
	acceptingCashoutsCrash = false
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
	s.repo.UpdateWinMultipliers(lastGameID, winMultiplier)
	s.repo.CreditingWinnings(lastGameID)
	s.StartPreparingCrash()
}

//
// FOR HANDLER
//

func (s *CrashService) GetAllRecords() ([]model.CrashRecord, error) {
	var allRecords []model.CrashRecord
	allRecords, err := s.repo.GetAllRecords()
	if err != nil {
		return allRecords, err
	}
	return allRecords, nil
}
