package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type DropPostgres struct {
	db *gorm.DB
}

func NewDropPostgres(db *gorm.DB) *DropPostgres {
	return &DropPostgres{db: db}
}

func (r *DropPostgres) NewDrop(itemId int) (model.Item, error) {
	r.db.Model(&model.DropRecord{}).Create(&model.DropRecord{})
	var item model.Item
	if err := r.db.Model(&model.Item{}).Where("id = ?", itemId).First(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *DropPostgres) GetLastDrops() ([]model.Item, error) {
	var drops []model.DropRecord
	var items []model.Item
	var item model.Item
	if err := r.db.Model(&model.DropRecord{}).Order("id desc").Limit(7).Find(&drops).Error; err != nil {
		return items, err
	}
	for _, drop := range drops {
		r.db.Model(&model.Item{}).Where("id = ?", drop.ItemID).First(&item)
		items = append(items, item)
	}
	return items, nil
}
