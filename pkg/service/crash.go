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
	repo         repository.Crash
	fakeBetsRepo repository.FakeBets
}

func NewCrashService(repo repository.Crash, fakeBetsRepo repository.FakeBets) *CrashService {
	return &CrashService{repo: repo, fakeBetsRepo: fakeBetsRepo}
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
	GameId         int     `json:"game_id"`
	PlayerID       string  `json:"player_id"`
	PlayerNickname string  `json:"player_nickname"`
	Multiplier     float64 `json:"multiplier"`
}

type ResponseCrash struct {
	GameID          int     `json:"game_id"`
	Status          string  `json:"status"`
	Multiplier      float64 `json:"multiplier"`
	TimeBeforeStart float64 `json:"time_before_start"`
	Length          float64 `json:"length"`
}

type InfoAboutCrashBet struct {
	PlayerID       string  `json:"player_id"`
	PlayerNickname string  `json:"player_nickname"`
	PlayerPhoto    int     `json:"player_photo"`
	Amount         float64 `json:"amount"`
	UserMultiplier float64 `json:"user_multiplier"`
	Winning        float64 `json:"winning"`
}

type BetsAtLastCrashGame struct {
	Bets []InfoAboutCrashBet `json:"bets"`
}

type PreparingCrashData struct {
	Status           string    `json:"status"`
	NewGameStartTime time.Time `json:"new_game_start_time"`
}

var startCrash = false
var betsAtLastCrashGame BetsAtLastCrashGame
var responseCrash = ResponseCrash{0, "Crashed", 1.0, 10.0, 0.0}
var clientsCrash = make(map[*ClientCrash]bool)
var clientsMutexCrash = &sync.Mutex{}
var winMultiplier = 0.0
var u = 0.0
var deltaCrash = 0.0
var lastCrashGameID int
var newGameStartTime time.Time

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
			fmt.Println("Read error:", err)
			break
		}
		if acceptingBetsCrash {
			var bet BetMessageCrash
			if err = json.Unmarshal(message, &bet); err != nil {
				fmt.Println("Invalid bet format:", err)
				continue
			}
			newBet := model.BetCrash{
				GameId: bet.GameId,
				UserID: bet.PlayerID,
				Amount: bet.Amount,
			}
			errorStr := s.repo.NewBetCrash(newBet)
			go s.AddBetCrashToResponse(bet.PlayerID, bet.Amount)
			fmt.Println(errorStr)
		} else if acceptingCashoutsCrash {
			var cashout CashoutMessageCrash
			if err = json.Unmarshal(message, &cashout); err != nil {
				fmt.Println("Invalid bet format:", err)
				continue
			}
			if cashout.PlayerID == "fake" {
				go s.CashOutFakeBets(cashout.PlayerNickname, cashout.Multiplier)
			} else if cashout.GameId == lastCrashGameID {
				errorStr := s.repo.NewCashoutCrash(cashout.GameId, cashout.PlayerID, cashout.Multiplier)
				go s.UpdateSavedBetCrash(cashout.PlayerID, cashout.Multiplier)
				fmt.Println(errorStr)
			}
		}
	}

	clientsMutexCrash.Lock()
	delete(clientsCrash, client)
	clientsMutexCrash.Unlock()
}

func (s *CrashService) CheckStatusOfStartCrash() {
	if startCrash == false {
		responseCrash.Status = "Stopped"
		clientsMutexCrash.Lock()
		for client := range clientsCrash {
			err := client.conn.WriteJSON(responseCrash.Status)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsCrash, client)
			}
		}
		clientsMutexCrash.Unlock()
		time.Sleep(1 * time.Second)
		s.CheckStatusOfStartCrash()
	} else {
		s.StartPreparingCrash()
	}
}

func (s *CrashService) ChangeStatusOfStartCrash(statusFromFront bool) {
	startCrash = statusFromFront
}

func (s *CrashService) CrashGame() {
	s.CheckStatusOfStartCrash()
}

func (s *CrashService) StartPreparingCrash() {
	betsAtLastCrashGame = BetsAtLastCrashGame{}
	clientsMutexCrash.Lock()
	for client := range clientsCrash {
		err := client.conn.WriteJSON(betsAtLastCrashGame)
		if err != nil {
			log.Println("Write error:", err)
			client.conn.Close()
			delete(clientsCrash, client)
		}
	}
	clientsMutexCrash.Unlock()
	acceptingBetsCrash = true
	responseCrash.Length = 0.0
	responseCrash.Status = "Pending"
	u = rand.Float64()
	winMultiplier = math.Pow(1-u, -1/2.25)
	//winMultiplier = 1.01
	lastGame, err := s.repo.GetLastCrashRecord()
	if err != nil {
		log.Fatal(err)
	}
	lastCrashGameID = lastGame.ID + 1
	responseCrash.GameID = lastCrashGameID
	s.PreparingCrash()
}

func (s *CrashService) PreparingCrash() {
	go s.GenerateFakeBetsCrash()
	newGameStartTime = time.Now().Add(10 * time.Second)
	clientsMutexCrash.Lock()
	for client := range clientsCrash {
		err := client.conn.WriteJSON(PreparingCrashData{
			Status:           "Penging",
			NewGameStartTime: newGameStartTime,
		})
		if err != nil {
			log.Println("Write error:", err)
			client.conn.Close()
			delete(clientsCrash, client)
		}
	}
	clientsMutexCrash.Unlock()
	time.Sleep(10 * time.Second)
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
		if responseCrash.Length <= 100.0 {
			responseCrash.Length += 0.4
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
	go s.repo.NewCrashRecord(winMultiplier)
	s.EndCrash()
}

func (s *CrashService) EndCrash() {
	acceptingCashoutsCrash = false
	responseCrash.Status = "Crashed"
	wayTo150 := (200 - responseCrash.Length) / 300
	for time_before_pending := 300; time_before_pending >= 0; time_before_pending-- {
		time.Sleep(10 * time.Millisecond)
		responseCrash.Length += wayTo150
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
	s.repo.UpdateWinMultipliers(lastCrashGameID, winMultiplier)
	s.repo.CreditingWinningsCrash(lastCrashGameID)
	if startCrash {
		s.StartPreparingCrash()
	} else {
		s.CheckStatusOfStartCrash()
	}
}

func (s *CrashService) AddBetCrashToResponse(userId string, amount float64) {
	nickname, photo, err := s.repo.GetUsersPhotoAndNickForCrash(userId)
	if err != nil {
		return
	}
	infoAboutCrashBet := InfoAboutCrashBet{
		PlayerID:       userId,
		PlayerNickname: nickname,
		Amount:         amount,
		UserMultiplier: 0,
		Winning:        0,
		PlayerPhoto:    photo,
	}
	betsAtLastCrashGame.Bets = append(
		betsAtLastCrashGame.Bets,
		infoAboutCrashBet,
	)
	clientsMutexCrash.Lock()
	for client := range clientsCrash {
		err := client.conn.WriteJSON(infoAboutCrashBet)
		if err != nil {
			log.Println("Write error:", err)
			client.conn.Close()
			delete(clientsCrash, client)
		}
	}
	clientsMutexCrash.Unlock()
}

func (s *CrashService) UpdateSavedBetCrash(userId string, multiplier float64) {
	for betsInCurrentGame := range betsAtLastCrashGame.Bets {
		if betsAtLastCrashGame.Bets[betsInCurrentGame].PlayerID == userId {
			currentWinning := math.Round(betsAtLastCrashGame.Bets[betsInCurrentGame].Amount*multiplier*100.0) / 100.0
			currentMultiplier := math.Round(multiplier*100.0) / 100.0
			betsAtLastCrashGame.Bets[betsInCurrentGame].UserMultiplier = currentMultiplier
			betsAtLastCrashGame.Bets[betsInCurrentGame].Winning = currentWinning
			clientsMutexCrash.Lock()
			for client := range clientsCrash {
				err := client.conn.WriteJSON(betsAtLastCrashGame.Bets[betsInCurrentGame])
				if err != nil {
					log.Println("Write error:", err)
					client.conn.Close()
					delete(clientsCrash, client)
				}
			}
			clientsMutexCrash.Unlock()
			break
		}
	}
}

func getRandomElementsForCrash(arr []model.FakeBets) []model.FakeBets {
	// Получаем длину массива
	length := len(arr)

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Создаем слайс для результата
	var result []model.FakeBets

	// Проверяем, что массив не пустой
	if length == 0 {
		return result
	}

	// Перемешиваем исходный массив
	shuffled := make([]model.FakeBets, length)
	copy(shuffled, arr)
	rand.Shuffle(length, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Определяем количество случайных элементов (до половины длины массива)
	randomCount := rand.Intn(length/2) + 1

	// Добавляем случайные элементы в результат
	for i := 0; i < randomCount; i++ {
		result = append(result, shuffled[i])
	}

	return result
}

func randomIntCrash(min, max int) float64 {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Генерируем случайное целое число в диапазоне [min, max]
	randomInt := rand.Intn(max-min+1) + min

	// Преобразуем его в формат float64 с двумя нулями
	return float64(randomInt)
}

func (s *CrashService) GenerateFakeBetsCrash() {
	users, err := s.fakeBetsRepo.GetFakeUsers()
	if err != nil {
		return
	}
	fakeBets := getRandomElementsForCrash(users)
	maxDelay := 10000 / len(fakeBets)
	var delay int
	var infoAboutFakeCrashBet InfoAboutCrashBet
	for _, fakeBet := range fakeBets {
		infoAboutFakeCrashBet = InfoAboutCrashBet{
			PlayerID:       "fake",
			PlayerNickname: fakeBet.Name,
			Amount:         randomIntCrash(10, 500),
			UserMultiplier: 0,
			Winning:        0,
			PlayerPhoto:    fakeBet.Photo,
		}
		betsAtLastCrashGame.Bets = append(
			betsAtLastCrashGame.Bets,
			infoAboutFakeCrashBet,
		)
		clientsMutexCrash.Lock()
		for client := range clientsCrash {
			err := client.conn.WriteJSON(infoAboutFakeCrashBet)
			if err != nil {
				log.Println("Write error:", err)
				client.conn.Close()
				delete(clientsCrash, client)
			}
		}
		clientsMutexCrash.Unlock()
		delay = int(randomIntCrash(0, maxDelay))
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

func (s *CrashService) CashOutFakeBets(name string, multiplier float64) {
	for betsInCurrentGame := range betsAtLastCrashGame.Bets {
		if betsAtLastCrashGame.Bets[betsInCurrentGame].PlayerNickname == name {
			currentWinning := math.Round(betsAtLastCrashGame.Bets[betsInCurrentGame].Amount*multiplier*100.0) / 100.0
			currentMultiplier := math.Round(multiplier*100.0) / 100.0
			betsAtLastCrashGame.Bets[betsInCurrentGame].UserMultiplier = currentMultiplier
			betsAtLastCrashGame.Bets[betsInCurrentGame].Winning = currentWinning
			clientsMutexCrash.Lock()
			for client := range clientsCrash {
				err := client.conn.WriteJSON(betsAtLastCrashGame.Bets[betsInCurrentGame])
				if err != nil {
					log.Println("Write error:", err)
					client.conn.Close()
					delete(clientsCrash, client)
				}
			}
			clientsMutexCrash.Unlock()
			break
		}
	}
}

//
// FOR HANDLER
//

func (s *CrashService) GetAllRecords() ([]model.CrashRecord, error) {
	var allRecords []model.CrashRecord
	allRecords, err := s.repo.GetAllCrashRecords()
	if err != nil {
		return allRecords, err
	}
	return allRecords, nil
}

func (s *CrashService) InitCrashBetsForNewClient() BetsAtLastCrashGame {
	return betsAtLastCrashGame
}
