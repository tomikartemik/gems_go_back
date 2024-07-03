package repository

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
	"strconv"
)

type ReplenishmentPostgres struct {
	db *gorm.DB
}

func NewReplenishmentPostgres(db *gorm.DB) *ReplenishmentPostgres {
	return &ReplenishmentPostgres{db: db}
}

func (r *ReplenishmentPostgres) NewReplenishment(userID string, amount float64) (string, string, error) {
	newReplenishment := model.Replenishment{
		UserID:  userID,
		Amount:  amount,
		Success: false,
	}
	result := r.db.Model(&model.Replenishment{}).Create(&newReplenishment)
	if result.Error != nil {
		return "", "", result.Error
	}
	var user model.User
	err := r.db.Model(&model.User{}).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return "", "", err
	}
	fmt.Println("repo: ", strconv.Itoa(newReplenishment.ID), user.Email)
	return strconv.Itoa(newReplenishment.ID), user.Email, nil
}
