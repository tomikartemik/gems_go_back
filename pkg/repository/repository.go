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
	Replenishment
}

type User interface {
	CreateUser(user model.User) (schema.ShowUser, error)
	SignIn(mail, password string) (schema.ShowUser, error)
	UpdateUser(id string, user schema.InputUser) (schema.ShowUser, error)
	GetUserInventory(userId string) ([]model.ItemOfInventory, error)
	GetUserById(id string) (schema.ShowUser, error)
	//AddItemToInventory(userId string, itemId int) (model.UserItem, error)
	SellAnItem(userId string, user_item_id int) error
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
	NewCaseRecord(caseId int) error
	GetAllCaseRecords() ([]schema.CaseInfo, error)
	AddItemToInventoryAndChangeBalance(userId string, itemId int) error
}

type Crash interface {
	NewCrashRecord(winMultiplier float64) error
	GetAllCrashRecords() ([]model.CrashRecord, error)
	GetLastCrashRecord() (model.CrashRecord, error)
	NewBetCrash(newBet model.BetCrash) string
	NewCashoutCrash(gameID int, userID string, userMultiplier float64) string
	UpdateWinMultipliers(gameID int, winMultiplier float64) string
	CreditingWinningsCrash(gameID int) string
	GetUsersPhotoAndNickForCrash(userId string) (string, error)
}

type Roulette interface {
	NewRouletteRecord(winCell int) error
	GetAllRouletteRecords() ([]model.RouletteRecord, error)
	GetLastRouletteRecord() (model.RouletteRecord, error)
	NewBetRoulette(newBet model.BetRoulette) string
	UpdateWinCells(gameID int, winCell int) string
	CreditingWinningsRoulette(gameID int) string
	GetUsersPhotoAndNickForRoulette(userId string) (string, error)
}

type Replenishment interface {
	NewReplenishment(userID string, amount float64) (string, string, error)
	AcceptReplenishment(replenishmentID int) error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:          NewUserPostgres(db),
		Item:          NewItemPostgres(db),
		Case:          NewCasePostgres(db),
		Crash:         NewCrashPostgres(db),
		Roulette:      NewRoulettePostgres(db),
		Replenishment: NewReplenishmentPostgres(db),
	}
}
