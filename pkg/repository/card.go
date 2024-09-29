package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type CardPostgres struct {
	db *gorm.DB
}

func NewCardPostgres(db *gorm.DB) *CardPostgres {
	return &CardPostgres{db: db}
}

func (r *CardPostgres) UpdateCard(card model.Card) error {
	err := r.db.Model(&model.Card{}).Where("type = main").Update("name", card.Name).Update("number", card.Number).Update("bank", card.Bank).Error
	return err
}

func (r *CardPostgres) GetCard() (model.Card, error) {
	var card model.Card
	err := r.db.First(&card).Error
	if err != nil {
		return card, err
	}
	return card, nil
}
