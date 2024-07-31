package service

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"gems_go_back/pkg/schema"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/websocket"
)

type User interface {
	CreateUser(user model.User) (schema.ShowUser, error)
	ParseToken(token string) (string, error)
	UpdateUser(id string, user schema.InputUser) (schema.ShowUser, error)
	GetUserById(id string) (schema.UserWithItems, error)
	SignIn(email string, password string) (SignInResponse, error)
	SellItem(userId string, userItemId int) (schema.UserWithItems, error)
	SellAllItems(userId string) error
	ChangeAvatar(userId string) (int, error)
}

type Item interface {
	CreateItem(item model.Item) (model.ItemWithID, error)
	GetItem(id int) (model.ItemWithID, error)
	GetAllItems() ([]model.Item, error)
	UpdateItem(item model.ItemWithID) (model.ItemWithID, error)
}

type Case interface {
	CreateCase(newCase model.Case) (schema.ShowCase, error)
	GetCase(id int) (schema.ShowCase, error)
	AddItemsToCase(id int, caseItems []model.CaseItem) (schema.ShowCase, error)
	GetCaseItems(caseId int) ([]model.ItemWithID, error)
	GetAllCases() ([]schema.CaseInfo, error)
	UpdateCase(id int, updates model.Case) (schema.ShowCase, error)
	DeleteCase(caseId int) error
	OpenCase(userId string, caseId int) (model.ItemWithID, int, error)
}

type Crash interface {
	EditConnsCrash(ws *websocket.Conn)
	CrashGame()
	CheckStatusOfStartCrash()
	ChangeStatusOfStartCrash(start bool)
	StartPreparingCrash()
	PreparingCrash()
	StartGameCrash()
	GameCrash()
	EndCrash()
	GetAllRecords() ([]model.CrashRecord, error)
	UpdateSavedBetCrash(userId string, multiplier float64)
	AddBetCrashToResponse(userId string, amount float64)
	InitCrashBetsForNewClient() BetsAtLastCrashGame
}

type Roulette interface {
	EidtConnsRoulette(conn *websocket.Conn)
	RouletteGame()
	CheckStatusOfStartRoulette()
	ChangeStatusOfStartRoulette(statusFromFront bool)
	StartPreparingRoulette()
	PreparingRoulette()
	StartGameRoulette()
	GameRoulette()
	EndRoulette()
	GetAllRouletteRecords() ([]model.RouletteRecord, error)
	InitRouletteBetsForNewClient() BetsAtLastRouletteGame
}

type Replenishment interface {
	NewReplenishment(userId string, amount float64) (string, error)
	AcceptReplenishment(replenishmentID int)
}

type Withdraw interface {
	TelegramBot()
	CreateWithdraw(newWithdraw model.Withdraw) error
	HandleUpdatesTelegram(bot *tgbotapi.BotAPI)
	GetUsersWithdraws(userId string) ([]model.Withdraw, error)
	GetPositionPrices() []model.Price
}

type Online interface {
	GetOnline() int
	SetOnline()
}

type Drop interface {
	GetLastDrops() []DropResponse
	EditConnsDrop(conn *websocket.Conn)
	NewDrop(itemId int, dirty bool)
	DirtyMoves()
	DropWS()
}

type Service struct {
	User
	Item
	Case
	Crash
	Roulette
	Replenishment
	Withdraw
	Online
	Drop
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:          NewAuthService(repos.User),
		Item:          NewItemService(repos.Item),
		Case:          NewCaseService(repos.Case),
		Crash:         NewCrashService(repos.Crash),
		Roulette:      NewRouletteService(repos.Roulette),
		Replenishment: NewReplenishmentService(repos.Replenishment),
		Withdraw:      NewWithdrawService(repos.Withdraw),
		Online:        NewOnlineService(repos.Online),
		Drop:          NewDropService(repos.Drop),
	}
}
