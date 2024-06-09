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

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		NewUserPostgres(db),
		NewItemPostgres(db),
		NewCasePostgres(db),
	}
}
