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
	repo         repository.Roulette
	fakeBetsRepo repository.FakeBets
}

func NewRouletteService(repo repository.Roulette, fakeBetsRepo repository.FakeBets) *RouletteService {
	return &RouletteService{repo: repo, fakeBetsRepo: fakeBetsRepo}
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
	PlayerNickname string  `json:"player_nickname"`
	Image          int     `json:"image"`
	Amount         float64 `json:"amount"`
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

var startRoulette = false
var betsAtLastRouletteGame BetsAtLastRouletteGame
var clientsRoulette = make(map[*ClientRoulette]bool)
var clientsMutexRoulette = &sync.Mutex{}
var responseRoulette = ResponseRouletteStatus{0, "Pending", 0, 0.0}
var cells = []Cell{
	{2, 50},
	{3, 33},
	{5, 20},
	{10, 10},
	{100, 1},
}
var totalWeightInRoulette = 114
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

func (s *RouletteService) RouletteGame() {
	//s.repo.NewRouletteRecord(100)
	s.StartPreparingRoulette()
}

func (s *RouletteService) CheckStatusOfStartRoulette() {
	if startRoulette == false {
		responseRoulette.Status = "Stopped"
		clientsMutexRoulette.Lock()
		for client := range clientsRoulette {
			err := client.conn.WriteJSON(responseRoulette.Status)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsRoulette, client)
			}
		}
		clientsMutexRoulette.Unlock()
		time.Sleep(1 * time.Second)
		s.CheckStatusOfStartRoulette()
	} else {
		s.StartPreparingRoulette()
	}
}

func (s *RouletteService) ChangeStatusOfStartRoulette(statusFromFront bool) {
	startRoulette = statusFromFront
}

func (s *RouletteService) StartPreparingRoulette() {
	betsAtLastRouletteGame = BetsAtLastRouletteGame{}
	clientsMutexRoulette.Lock()
	for client := range clientsRoulette {
		err := client.conn.WriteJSON(betsAtLastRouletteGame)
		if err != nil {
			log.Println("Write error:", err)
			client.conn.Close()
			delete(clientsRoulette, client)
		}
	}
	clientsMutexRoulette.Unlock()
	betsAtLastRouletteGame.MainAmount = 0.0
	acceptingBetsRoulette = true
	responseRoulette.Cell = 0
	responseRoulette.Status = "Pending"
	lastGame, err := s.repo.GetLastRouletteRecord()
	if err != nil {
		fmt.Println(err)
	}
	lsatRouletteGameID = lastGame.ID + 1
	responseRoulette.GameID = lsatRouletteGameID
	s.PreparingRoulette()
}

func (s *RouletteService) PreparingRoulette() {
	go s.GenerateFakeBetsRoulette()
	for time_before_start := 1000.0; time_before_start >= 0; time_before_start-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutexRoulette.Lock()
		responseRoulette.TimeBeforeStart = time_before_start / 100.0
		for client := range clientsRoulette {
			err := client.conn.WriteJSON(responseRoulette)
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
	responseRoulette.Status = "Playing"
	acceptingBetsRoulette = false
	randomNumber := rand.Intn(totalWeightInRoulette)

	for _, choosenCell := range cells {
		if randomNumber < choosenCell.Weight {
			winCell = choosenCell.Value
			break
		}
		randomNumber -= choosenCell.Weight
	}
	responseRoulette.Cell = winCell
	s.GameRoulette()
}

func (s *RouletteService) GameRoulette() {
	for time_before_end := 700.0; time_before_end >= 0; time_before_end-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutexRoulette.Lock()
		for client := range clientsRoulette {
			err := client.conn.WriteJSON(responseRoulette)
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
	responseRoulette.Status = "End"
	for time_before_pending := 300; time_before_pending >= 0; time_before_pending-- {
		time.Sleep(10 * time.Millisecond)
		clientsMutexRoulette.Lock()
		for client := range clientsRoulette {
			err := client.conn.WriteJSON(responseRoulette)
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
	if startRoulette {
		s.StartPreparingRoulette()
	} else {
		s.CheckStatusOfStartRoulette()
	}
}

func (s *RouletteService) AddRouletteBetToResponse(userID string, amount float64, cell int) {
	var betMessageRouletteResponse BetMessageRouletteResponse
	betsAtLastRouletteGame.MainAmount += amount
	playerNickname, photo, err := s.repo.GetUsersPhotoAndNickForRoulette(userID)
	if err != nil {
		log.Println(err)
	}
	betMessageRouletteResponse.PlayerNickname = playerNickname
	betMessageRouletteResponse.Amount = amount
	betMessageRouletteResponse.Image = photo
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
	for client := range clientsRoulette {
		err = client.conn.WriteJSON(betsAtLastRouletteGame)
		if err != nil {
			log.Println("Write error:", err)
			client.conn.Close()
			delete(clientsRoulette, client)
		}
	}
	clientsMutexRoulette.Unlock()
}

func getRandomElementsForRoulette(arr []model.FakeBets) []model.FakeBets {
	// Получаем длину массива
	length := len(arr) / 3

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Выбираем случайное число от 0 до длины массива
	randomCount := rand.Intn(length) + 1

	// Создаем слайс для результата
	var result []model.FakeBets

	// Добавляем случайные элементы в результат
	for i := 0; i < randomCount; i++ {
		// Выбираем случайный индекс и добавляем элемент в результат
		randomIndex := rand.Intn(length)
		result = append(result, arr[randomIndex])
	}

	return result
}

func randomIntRoulette(min, max int) float64 {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Генерируем случайное целое число в диапазоне [min, max]
	randomInt := rand.Intn(max-min+1) + min

	// Преобразуем его в формат float64 с двумя нулями
	return float64(randomInt)
}

func getRandomCell() int {
	options := []int{2, 3, 5, 10, 100}
	rand.Seed(time.Now().UnixNano())        // Инициализация генератора случайных чисел
	return options[rand.Intn(len(options))] // Выбор случайного числа из options
}

func (s *RouletteService) GenerateFakeBetsRoulette() {
	users, err := s.fakeBetsRepo.GetFakeUsers()
	if err != nil {
		return
	}
	fakeBets := getRandomElementsForRoulette(users)
	maxDelay := 10000 / len(fakeBets)
	var delay int
	var cell int
	var amount float64
	var infoAboutFakeRouletteBet BetMessageRouletteResponse
	for _, fakeBet := range fakeBets {
		amount = randomIntRoulette(10, 500)
		infoAboutFakeRouletteBet = BetMessageRouletteResponse{
			PlayerNickname: fakeBet.Name,
			Amount:         randomIntRoulette(10, 500),
			Image:          fakeBet.Photo,
		}

		cell = getRandomCell()
		betsAtLastRouletteGame.MainAmount += amount

		if cell == 2 {
			betsAtLastRouletteGame.Bet2 = append(betsAtLastRouletteGame.Bet2, infoAboutFakeRouletteBet)
		} else if cell == 3 {
			betsAtLastRouletteGame.Bet3 = append(betsAtLastRouletteGame.Bet3, infoAboutFakeRouletteBet)
		} else if cell == 5 {
			betsAtLastRouletteGame.Bet5 = append(betsAtLastRouletteGame.Bet5, infoAboutFakeRouletteBet)
		} else if cell == 10 {
			betsAtLastRouletteGame.Bet10 = append(betsAtLastRouletteGame.Bet10, infoAboutFakeRouletteBet)
		} else if cell == 100 {
			betsAtLastRouletteGame.Bet100 = append(betsAtLastRouletteGame.Bet100, infoAboutFakeRouletteBet)
		}
		clientsMutexRoulette.Lock()
		for client := range clientsRoulette {
			err = client.conn.WriteJSON(betsAtLastRouletteGame)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsRoulette, client)
			}
		}
		clientsMutexRoulette.Unlock()

		delay = int(randomIntCrash(0, maxDelay))
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

func (s *RouletteService) GetAllRouletteRecords() ([]model.RouletteRecord, error) {
	var lastRecords []model.RouletteRecord
	lastRecords, err := s.repo.GetAllRouletteRecords()
	if err != nil {
		return lastRecords, err
	}
	return lastRecords, nil
}

func (s *RouletteService) InitRouletteBetsForNewClient() BetsAtLastRouletteGame {
	return betsAtLastRouletteGame
}
