package repository

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type DropPostgres struct {
	db *gorm.DB
}

func NewDropPostgres(db *gorm.DB) *DropPostgres {
	return &DropPostgres{db: db}
}

func (r *DropPostgres) NewDrop(itemId int, dirty bool) (model.Item, error) {
	r.db.Model(&model.DropRecord{}).Create(&model.DropRecord{ItemID: itemId, Dirty: dirty})
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
	fmt.Println(drops)
	for _, drop := range drops {
		r.db.Model(&model.Item{}).Where("id = ?", drop.ItemID).First(&item)
		items = append(items, item)
	}
	fmt.Println(items)
	return items, nil
}

func (r *DropPostgres) GetItemsIds() ([]int, error) {
	var ids []int
	var items []model.Item
	if err := r.db.Model(&model.Item{}).Find(&items).Error; err != nil {
		return ids, err
	}
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids, nil
}
