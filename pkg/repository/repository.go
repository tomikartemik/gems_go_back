package repository

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/schema"
	"gorm.io/gorm"
)

type Repository struct {
	User
	Item
	Case
	Crash
	Roulette
}

type User interface {
	CreateUser(user model.User) (schema.ShowUser, error)
	SignIn(mail, password string) (schema.ShowUser, error)
	UpdateUser(id string, user schema.InputUser) (schema.ShowUser, error)
	GetUserInventory(userId string) ([]model.ItemWithID, error)
	GetUserById(id string) (schema.ShowUser, error)
	AddItemToInventory(userId string, itemId int) (model.UserItem, error)
}

type Item interface {
	CreateItem(item model.Item) (model.ItemWithID, error)
	GetItem(id int) (model.ItemWithID, error)
	GetAllItems() ([]model.Item, error)
	UpdateItem(item model.ItemWithID) (model.ItemWithID, error)
}

type Case interface {
	CreateCase(newCase schema.CaseInput) (int, error)
	GetCaseInfo(id int) (model.Case, error)
	AddItemsToCase(id int, caseItems []model.CaseItem) (model.Case, error)
	DeleteItemsFromCase(id int) error
	GetCaseItems(caseId int) ([]model.ItemWithID, error)
	GetAllCases() ([]schema.CaseInfo, error)
	UpdateCase(id int, newCase schema.CaseInput) (schema.ShowCase, error)
	DeleteCase(id int) error
	GetItemsWithWeights(id int) ([]model.CaseItem, error)
	GetChosenItem(id int) (model.ItemWithID, error)
}

type Crash interface {
	NewRecord(winMultiplier float64) error
	GetAllRecords() ([]model.CrashRecord, error)
	GetLastRecord() (model.CrashRecord, error)
	NewBetCrash(newBet model.BetCrash) string
	NewCashoutCrash(gameID int, userID string, userMultiplier float64) string
	UpdateWinMultipliers(gameID int, winMultiplier float64) string
	CreditingWinnings(gameID int) string
}

type Roulette interface {
	EchoMSGRoulette(msg string) string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:     NewUserPostgres(db),
		Item:     NewItemPostgres(db),
		Case:     NewCasePostgres(db),
		Crash:    NewCrashPostgres(db),
		Roulette: NewRoulettePostgres(db),
	}
}
