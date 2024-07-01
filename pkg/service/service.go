package service

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"gems_go_back/pkg/schema"
	"github.com/gorilla/websocket"
)

type User interface {
	CreateUser(user model.User) (schema.ShowUser, error)
	GenerateToken(mail, password string) (string, error)
	ParseToken(token string) (string, error)
	UpdateUser(id string, user schema.InputUser) (schema.ShowUser, error)
	GetUserById(id string) (schema.UserWithItems, error)
	AddItemToInventory(userId string, itemId int) (schema.ShowUser, error)
	SellAnItem(userId string, userItemId int) (schema.UserWithItems, error)
	SignIn(email string, password string) (schema.UserWithItems, error)
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
	OpenCase(caseId int) (model.ItemWithID, error)
	GetAllCaseRecords() ([]schema.CaseInfo, error)
}

type Crash interface {
	EditConnsCrash(ws *websocket.Conn)
	BroadcastTimeCrash()
	StartPreparingCrash()
	PreparingCrash()
	StartGameCrash()
	GameCrash()
	EndCrash()
	GetAllRecords() ([]model.CrashRecord, error)
}

type Roulette interface {
	EidtConnsRoulette(conn *websocket.Conn)
	BroadcastTimeRoulette()
	StartPreparingRoulette()
	PreparingRoulette()
	StartGameRoulette()
	GameRoulette()
	EndRoulette()
	GetAllRouletteRecords() ([]model.RouletteRecord, error)
}

type Service struct {
	User
	Item
	Case
	Crash
	Roulette
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:     NewAuthService(repos.User),
		Item:     NewItemService(repos.Item),
		Case:     NewCaseService(repos.Case),
		Crash:    NewCrashService(repos.Crash),
		Roulette: NewRouletteService(repos.Roulette),
	}
}
