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

func (r *OwnReplenishmentPostgres) GetReplenishments() ([]model.OwnReplenishment, error) {
	replenishments := []model.OwnReplenishment{}
	err := r.db.Find(&replenishments).Where("status = Processing").Error
	if err != nil {
		return nil, err
	}
	return replenishments, nil
}
