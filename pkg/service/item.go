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
