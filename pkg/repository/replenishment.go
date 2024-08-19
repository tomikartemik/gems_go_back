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

func (r *ReplenishmentPostgres) GetReward(promo, userID string) float64 {
	reward := 1.0
	promoInfo := &model.Promo{}
	if err := r.db.Model(&model.Promo{}).Where("promo = ?", promo).First(&promoInfo).Error; err != nil {
		return reward
	}
	if err := r.db.Model(&model.PromoUsage{}).Where("promo_id = ?", promoInfo.ID).Where("user_id = ?", userID).Error; err == gorm.ErrRecordNotFound {
		r.db.Model(&model.PromoUsage{}).Create(&model.PromoUsage{
			PromoID: promoInfo.ID,
			UserID:  userID,
		})
		return promoInfo.Reward
	}
	return reward
}
