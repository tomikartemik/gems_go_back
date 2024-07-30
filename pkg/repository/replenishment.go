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
	var user model.User
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return "", "", err
	}

	newReplenishment := model.Replenishment{
		UserID:  userID,
		Amount:  amount,
		Success: false,
	}
	result := r.db.Model(&model.Replenishment{}).Create(&newReplenishment)
	if result.Error != nil {
		return "", "", result.Error
	}
	fmt.Println("repo: ", strconv.Itoa(newReplenishment.ID), user.Email)
	return strconv.Itoa(newReplenishment.ID), user.Email, nil
}

func (r *ReplenishmentPostgres) AcceptReplenishment(replenishmentID int) error {
	var replenishment model.Replenishment
	err := r.db.Model(&model.Replenishment{}).Where("id = ?", replenishmentID).First(&replenishment).Error
	if err != nil {
		return err
	}
	err = r.db.Model(&model.Replenishment{}).Where("id = ?", replenishmentID).Update("success", true).Error
	if err != nil {
		return err
	}
	err = r.db.Model(&model.User{}).Where("id = ?", replenishment.UserID).Update("balance", gorm.Expr("balance + ?", replenishment.Amount)).Error
	if err != nil {
		return err
	}
	return nil
}
