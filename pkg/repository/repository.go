package repository

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/schema"
	"gorm.io/gorm"
	"mime/multipart"
)

type Repository struct {
	User
	Item
	Case
	Crash
	Roulette
	Replenishment
	Withdraw
	Online
	Drop
	Receipt
	Admin
	OwnReplenishment
	Card
	FakeBets
}

type User interface {
	CreateUser(user model.User) (schema.ShowUser, error)
	SignIn(mail, password string) (schema.ShowUser, error)
	UpdateUser(id string, user schema.InputUser) (schema.ShowUser, error)
	GetUserInventory(userId string) ([]model.ItemOfInventory, error)
	GetUserById(id string) (schema.ShowUser, error)
	SellItem(userId string, user_item_id int) error
	SellAllItem(userId string) error
	ChangeAvatar(userId string, newPhoto int) error
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
	CheckThePossibilityOfPurchasing(userId string, caseId int) bool
	GetChosenItem(id int) (model.ItemWithID, error)
	AddItemToInventoryAndChangeBalance(userId string, itemId int, caseId int) (int, error)
}

type Crash interface {
	NewCrashRecord(winMultiplier float64) error
	GetAllCrashRecords() ([]model.CrashRecord, error)
	GetLastCrashRecord() (model.CrashRecord, error)
	NewBetCrash(newBet model.BetCrash) string
	NewCashoutCrash(gameID int, userID string, userMultiplier float64) string
	UpdateWinMultipliers(gameID int, winMultiplier float64) string
	CreditingWinningsCrash(gameID int) string
	GetUsersPhotoAndNickForCrash(userId string) (string, int, error)
}

type Roulette interface {
	NewRouletteRecord(winCell int) error
	GetAllRouletteRecords() ([]model.RouletteRecord, error)
	GetLastRouletteRecord() (model.RouletteRecord, error)
	NewBetRoulette(newBet model.BetRoulette) string
	UpdateWinCells(gameID int, winCell int) string
	CreditingWinningsRoulette(gameID int) string
	GetUsersPhotoAndNickForRoulette(userId string) (string, int, error)
}

type Replenishment interface {
	NewReplenishment(userID string, amount float64) (string, string, error)
	AcceptReplenishment(replenishmentID int) error
	GetReward(promo, userID string) (float64, error)
}

type Withdraw interface {
	CreateWithdraw(tx *gorm.DB, withdraw model.Withdraw) (model.Withdraw, error)
	BeginTransaction() *gorm.DB
	CompleteWithdraw(withdrawId int) error
	GetWithdraw(withdrawId int) (model.Withdraw, error)
	GetUsersWithdraws(userId string) ([]model.Withdraw, error)
	CancelWithdraw(withdrawId int) error
	ReturnMoneyBecauseCanceled(currentWithdraw model.Withdraw)
	GetPositionPrice(position string) (float64, error)
	GetPositionPrices() []model.Price
}

type Online interface {
	GetOnline() int
	SetOnline(usersOnline int)
}

type Drop interface {
	NewDrop(itemId int, dirty bool) (model.Item, error)
	GetItemsIds() ([]int, error)
}

type Receipt interface {
	UploadReceipt(file *multipart.FileHeader, filepath string) error
}

type Admin interface {
	CreateAdmin(admin model.Admin) error
	SignInAdmin(mail, password string) (model.Admin, error)
}

type OwnReplenishment interface {
	CreateReplenishment(replenishment model.OwnReplenishment) error
	GetReplenishments(sortOrder, status string) ([]model.OwnReplenishment, error)
	ChangeStatus(replenishmentID int, status string) (model.OwnReplenishment, error)
	ChangeBalance(userID string, amount float64) error
	GetLastId() (int, error)
}

type Card interface {
	UpdateCard(card model.Card) error
	GetCard() (model.Card, error)
}

type FakeBets interface {
	GetFakeUsers() ([]model.FakeBets, error)
	CreateFakeBet(user model.FakeBets)
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:             NewUserPostgres(db),
		Item:             NewItemPostgres(db),
		Case:             NewCasePostgres(db),
		Crash:            NewCrashPostgres(db),
		Roulette:         NewRoulettePostgres(db),
		Replenishment:    NewReplenishmentPostgres(db),
		Withdraw:         NewWithdrawPostgres(db),
		Online:           NewOnlinePostgres(db),
		Drop:             NewDropPostgres(db),
		Receipt:          NewReceiptPostgres(db),
		Admin:            NewAdminPostgres(db),
		OwnReplenishment: NewOwnReplenishmentPostgres(db),
		Card:             NewCardPostgres(db),
		FakeBets:         NewFakeBetsPostgres(db),
	}
}
