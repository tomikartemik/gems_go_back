package repository

import (
	"errors"
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

func (r *OwnReplenishmentPostgres) ChangeStatus(replenishmentID int, status string) (model.OwnReplenishment, error) {
	replenishment := model.OwnReplenishment{}
	err := r.db.Find(&replenishment, replenishmentID).Error
	if err != nil {
		return replenishment, err
	}
	if replenishment.Status != "Processing" {
		return replenishment, errors.New("replenishment is not processing")
	}
	replenishment.Status = status
	err = r.db.Save(&replenishment).Error
	if err != nil {
		return replenishment, err
	}
	return replenishment, nil
}

func (r *OwnReplenishmentPostgres) ChangeBalance(userID string, amount float64) error {
	user := model.User{}
	err := r.db.First(&user, userID).Error
	if err != nil {
		return err
	}
	user.Balance = user.Balance + amount
	return r.db.Save(&user).Error
}
