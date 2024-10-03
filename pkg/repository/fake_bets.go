package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type FakeBetsPostgres struct {
	db *gorm.DB
}

func NewFakeBetsPostgres(db *gorm.DB) *FakeBetsPostgres {
	return &FakeBetsPostgres{db: db}
}

func (r *FakeBetsPostgres) GetFakeUsers() ([]model.FakeBets, error) {
	var users []model.FakeBets
	r.db.Model(&model.FakeBets{}).Find(&users)
	return users, nil
}
