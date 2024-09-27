package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type OwnReplenishmentPostgres struct {
	db *gorm.DB
}

func NewOwnReplenishmentPostgres(db *gorm.DB) *OwnReplenishmentPostgres {
	return &OwnReplenishmentPostgres{db: db}
}

func (r *OwnReplenishmentPostgres) CreateReplenishment(replenishment model.OwnReplenishment) error {
	return r.db.Create(&replenishment).Error
}
