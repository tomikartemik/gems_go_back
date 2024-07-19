package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type RoulettePostgres struct {
	db *gorm.DB
}

func NewRoulettePostgres(db *gorm.DB) *RoulettePostgres {
	return &RoulettePostgres{db: db}
}

func (r *RoulettePostgres) NewRouletteRecord(winCell int) error {
	record := model.RouletteRecord{
		WinCell: winCell,
	}
	result := r.db.Model(&model.RouletteRecord{}).Create(&record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RoulettePostgres) GetAllRouletteRecords() ([]model.RouletteRecord, error) {
	var all_records []model.RouletteRecord
	err := r.db.Model(&model.RouletteRecord{}).Order("id desc").Limit(10).Find(&all_records).Error
	if err != nil {
		return nil, err
	}
	return all_records, nil
}

func (r *RoulettePostgres) GetLastRouletteRecord() (model.RouletteRecord, error) {
	var last_record model.RouletteRecord
	err := r.db.Model(&model.RouletteRecord{}).Order("id desc").Last(&last_record).Error
	if err != nil {
		return last_record, err
	}
	return last_record, nil
}

func (r *RoulettePostgres) NewBetRoulette(newBet model.BetRoulette) string {
	var user model.User
	err := r.db.Model(&model.User{}).Where("id = ?", newBet.UserID).First(&user).Error
	if err != nil {
		return "User not found!"
	}
	if user.Balance < newBet.Amount {
		return "Not enough money!"
	}
	err = r.db.Model(&model.User{}).Where("id = ?", user.Id).Update("balance", (user.Balance - newBet.Amount)).Error
	if err != nil {
		return "Хуйня какая-то"
	}
	err = r.db.Model(&model.BetRoulette{}).Create(&newBet).Error
	if err != nil {
		return "Хуйня какая-то"
	}
	return "OK"
}

func (r *RoulettePostgres) UpdateWinCells(gameID int, winCell int) string {
	err := r.db.Model(&model.BetRoulette{}).Where("game_id = ?", gameID).Update("win_cell", winCell).Error
	if err != nil {
		return "Pizda!"
	}
	return "OK"
}

func (r *RoulettePostgres) CreditingWinningsRoulette(gameID int) string {
	var bets []model.BetRoulette
	if err := r.db.Model(&model.BetRoulette{}).Where("game_id = ?", gameID).Find(&bets).Error; err != nil {
		return "Pizda!"
	}
	for _, bet := range bets {
		if bet.UserCell == bet.WinCell {
			winAmount := bet.Amount * float64(bet.WinCell)
			if err := r.db.Model(&model.User{}).Where("id = ?", bet.UserID).Update("balance", gorm.Expr("balance + ?", winAmount)); err != nil {
				return "Pizda!"
			}
		}
	}
	return "OK"
}

func (r *RoulettePostgres) GetUsersPhotoAndNickForRoulette(userId string) (string, error) {
	var user model.User
	err := r.db.Model(&model.User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Username, err
}
