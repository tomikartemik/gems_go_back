package repository

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ReplenishmentPostgres struct {
	db *gorm.DB
}

func NewReplenishmentPostgres(db *gorm.DB) *ReplenishmentPostgres {
	return &ReplenishmentPostgres{db: db}
}

func (r *ReplenishmentPostgres) NewReplenishment(userID string, amount float64) (string, string, error) {
	orderID := int(time.Now().UnixNano())
	newReplenishment := model.Replenishment{
		ID:      orderID,
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
	return strconv.Itoa(orderID), user.Email, nil
}

func (r *ReplenishmentPostgres) AcceptReplenishment(replenishmentID int) error {
	var replenishment model.Replenishment
	err := r.db.Model(&model.Replenishment{}).Where("id = ?", replenishmentID).First(&replenishment).Error
	if err != nil {
		return err
	}
	if replenishment.Success == true {
		return nil
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

func (r *ReplenishmentPostgres) GetReward(promo, userID string) (float64, error) {
	defaultReward := 1.0
	promoInfo := &model.Promo{}
	tx := r.db.Begin()
	if tx.Error != nil {
		return defaultReward, tx.Error
	}
	if err := tx.Model(&model.Promo{}).Where("promo = ?", promo).First(&promoInfo).Error; err != nil {
		tx.Rollback()
		return defaultReward, err
	}

	usage := &model.PromoUsage{}
	if err := tx.Model(&model.PromoUsage{}).Where("promo_id = ?", promoInfo.ID).Where("user_id = ?", userID).First(&usage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newUsage := &model.PromoUsage{
				PromoID: promoInfo.ID,
				UserID:  userID,
			}
			if err := tx.Create(newUsage).Error; err != nil {
				tx.Rollback()
				return defaultReward, err
			}
			tx.Commit()
			return promoInfo.Reward, nil
		}
		tx.Rollback()
		return defaultReward, err
	}

	tx.Commit()
	return defaultReward, nil
}
