package service

import (
	"encoding/json"
	"fmt"
	"gems_go_back/pkg/model"
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

type BetMessageRoulette struct {
	GameId   int     `json:"game_id"`
	PlayerID string  `json:"player_id"`
	Amount   float64 `json:"amount"`
	Cell     int     `json:"cell"`
}

type BetMessageRouletteResponse struct {
	PlayerID       string `json:"player_id"`
	PlayerNickname string `json:"player_nickname"`
	//Image    string `json:"image"`
	Amount float64 `json:"amount"`
}

type BetsAtLastRouletteGame struct {
	MainAmount float64                      `json:"main_amount"`
	Bet2       []BetMessageRouletteResponse `json:"bet2"`
	Bet3       []BetMessageRouletteResponse `json:"bet3"`
	Bet5       []BetMessageRouletteResponse `json:"bet5"`
	Bet10      []BetMessageRouletteResponse `json:"bet10"`
	Bet100     []BetMessageRouletteResponse `json:"bet100"`
}

type ResponseRouletteStatus struct {
	GameID          int     `json:"game_id"`
	Status          string  `json:"status"`
	Cell            int     `json:"cell"`
	TimeBeforeStart float64 `json:"time_before_start"`
}

type Cell struct {
	Value  int
	Weight int
}

var betsAtLastRouletteGame BetsAtLastRouletteGame
var clientsRoulette = make(map[*ClientRoulette]bool)
var clientsMutexRoulette = &sync.Mutex{}
var responeRoulette = ResponseRouletteStatus{0, "Pending", 0, 0.0}
var cells = []Cell{
	{2, 50},
	{3, 33},
	{5, 20},
	{10, 10},
	{100, 1},
}
var totalWeight = 114
var winCell = 0

var lsatRouletteGameID int
var acceptingBetsRoulette = true

func (s *RouletteService) EidtConnsRoulette(conn *websocket.Conn) {

	defer conn.Close()

	client := &ClientRoulette{conn: conn}
	clientsMutexRoulette.Lock()
	clientsRoulette[client] = true
	clientsMutexRoulette.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		if acceptingBetsRoulette {
			var bet BetMessageRoulette
			if err = json.Unmarshal(message, &bet); err != nil {
				log.Println("Invalid bet format:", err)
				continue
			}
			newBet := model.BetRoulette{
				GameId:   bet.GameId,
				UserID:   bet.PlayerID,
				Amount:   bet.Amount,
				UserCell: bet.Cell,
			}
			errorStr := s.repo.NewBetRoulette(newBet)
			go s.AddRouletteBetToResponse(bet.PlayerID, bet.Amount, bet.Cell)
			fmt.Println(errorStr)
		}
	}

	clientsMutexRoulette.Lock()
	delete(clientsRoulette, client)
	clientsMutexRoulette.Unlock()
}

func (s *RouletteService) BroadcastTimeRoulette() {
	//s.repo.NewRouletteRecord(100)
	s.StartPreparingRoulette()
}

func (s *RouletteService) StartPreparingRoulette() {
	betsAtLastRouletteGame = BetsAtLastRouletteGame{}
	betsAtLastRouletteGame.MainAmount = 0.0
	acceptingBetsRoulette = true
	responeRoulette.Cell = 0
	responeRoulette.Status = "Pending"
	lastGame, err := s.repo.GetLastRouletteRecord()
	if err != nil {
		fmt.Println(err)
	}
	lsatRouletteGameID = lastGame.ID + 1
	responeRoulette.GameID = lsatRouletteGameID
	s.PreparingRoulette()
}

func (s *RouletteService) PreparingRoulette() {
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
	s.StartGameRoulette()
}

func (s *RouletteService) StartGameRoulette() {
	responeRoulette.Status = "Playing"
	acceptingBetsRoulette = false
	randomNumber := rand.Intn(totalWeight)

	for _, choosenCell := range cells {
		if randomNumber < choosenCell.Weight {
			winCell = choosenCell.Value
			break
		}
		randomNumber -= choosenCell.Weight
	}
	responeRoulette.Cell = winCell
	s.GameRoulette()
}

func (s *RouletteService) GameRoulette() {
	for time_before_end := 700.0; time_before_end >= 0; time_before_end-- {
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
	go s.repo.NewRouletteRecord(winCell)
	s.EndRoulette()
}

func (s *RouletteService) EndRoulette() {
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
	s.repo.UpdateWinCells(lsatRouletteGameID, winCell)
	s.repo.CreditingWinningsRoulette(lsatRouletteGameID)
	s.StartPreparingRoulette()
}

func (s *RouletteService) AddRouletteBetToResponse(userID string, amount float64, cell int) {
	var betMessageRouletteResponse BetMessageRouletteResponse
	betsAtLastRouletteGame.MainAmount += amount
	playerNickname, err := s.repo.GetUsersPhotoAndNickForRoulette(userID)
	if err != nil {
		log.Println(err)
	}
	betMessageRouletteResponse.PlayerNickname = playerNickname
	betMessageRouletteResponse.Amount = amount
	betMessageRouletteResponse.PlayerID = userID
	if cell == 2 {
		betsAtLastRouletteGame.Bet2 = append(betsAtLastRouletteGame.Bet2, betMessageRouletteResponse)
	} else if cell == 3 {
		betsAtLastRouletteGame.Bet3 = append(betsAtLastRouletteGame.Bet3, betMessageRouletteResponse)
	} else if cell == 5 {
		betsAtLastRouletteGame.Bet5 = append(betsAtLastRouletteGame.Bet5, betMessageRouletteResponse)
	} else if cell == 10 {
		betsAtLastRouletteGame.Bet10 = append(betsAtLastRouletteGame.Bet10, betMessageRouletteResponse)
	} else if cell == 100 {
		betsAtLastRouletteGame.Bet100 = append(betsAtLastRouletteGame.Bet100, betMessageRouletteResponse)
	}
	clientsMutexRoulette.Lock()
	for client := range clientsCrash {
		err = client.conn.WriteJSON(betsAtLastRouletteGame)
		if err != nil {
			log.Println("Write error:", err)
			client.conn.Close()
			delete(clientsCrash, client)
		}
	}
	clientsMutexRoulette.Unlock()
}

func (s *RouletteService) GetAllRouletteRecords() ([]model.RouletteRecord, error) {
	var lastRecords []model.RouletteRecord
	lastRecords, err := s.repo.GetAllRouletteRecords()
	if err != nil {
		return lastRecords, err
	}
	return lastRecords, nil
}
