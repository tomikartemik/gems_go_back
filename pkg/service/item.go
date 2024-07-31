package service

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
)

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(item model.Item) (model.ItemWithID, error) {
	if item.Rarity < 10 {
		item.Color = "#fd9d2d"
	} else if item.Rarity < 20 {
		item.Color = "#6cbf01"
	} else if item.Rarity < 40 {
		item.Color = "#00989E"
	} else if item.Rarity < 50 {
		item.Color = "#db2f4c"
	} else {
		item.Color = "#ffffff"
	}
	return s.repo.CreateItem(item)
}

func (s *ItemService) GetItem(id int) (model.ItemWithID, error) {

	return s.repo.GetItem(id)
}

func (s *ItemService) GetAllItems() ([]model.Item, error) {

	return s.repo.GetAllItems()
}

func (s *ItemService) UpdateItem(item model.ItemWithID) (model.ItemWithID, error) {
	return s.repo.UpdateItem(item)
}
