package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type ItemPostgres struct {
	db *gorm.DB
}

func NewItemPostgres(db *gorm.DB) *ItemPostgres {
	return &ItemPostgres{db: db}
}

func (r *ItemPostgres) CreateItem(item model.Item) (model.ItemWithID, error) {
	var newItem model.ItemWithID
	result := r.db.Model(&model.Item{}).Create(&item)
	if result.Error != nil {
		return newItem, result.Error
	}
	id := item.ID
	newItem, err := r.GetItem(id)
	if err != nil {
		return newItem, err
	}

	return newItem, nil
}

func (r *ItemPostgres) GetItem(id int) (model.ItemWithID, error) {
	var item model.ItemWithID
	result := r.db.Model(&model.Item{}).Where("id = ?", id).First(&item)
	if result.Error != nil {
		return item, result.Error
	}
	return item, nil
}

func (r *ItemPostgres) GetAllItems() ([]model.Item, error) {
	var items []model.Item
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ItemPostgres) UpdateItem(item model.ItemWithID) (model.ItemWithID, error) {
	id := item.ID
	if err := r.db.Model(&model.Item{}).Where("id = ?", id).Updates(&item).Error; err != nil {
		return model.ItemWithID{}, err
	}
	updatedItem, _ := r.GetItem(id)
	return updatedItem, nil
}
