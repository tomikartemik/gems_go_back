package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type CrashPostgres struct {
	db *gorm.DB
}

func NewCrashPostgres(db *gorm.DB) *CrashPostgres {
	return &CrashPostgres{db: db}
}

func (r *CrashPostgres) NewRecord(winMuliplier float64) error {
	record := model.CrashRecord{
		WinMultiplier: winMuliplier,
	}
	result := r.db.Model(&model.CrashRecord{}).Create(&record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CrashPostgres) GetAllCrashRecords() ([]model.CrashRecord, error) {
	var all_records []model.CrashRecord
	err := r.db.Model(&model.CrashRecord{}).Order("id desc").Limit(10).Find(&all_records).Error
	if err != nil {
		return nil, err
	}
	return all_records, nil
}

func (r *CrashPostgres) GetLastCrashRecord() (model.CrashRecord, error) {
	var last_record model.CrashRecord
	err := r.db.Model(&model.CrashRecord{}).Order("id desc").Last(&last_record).Error
	if err != nil {
		return last_record, err
	}
	return last_record, nil
}

func (r *CrashPostgres) NewBetCrash(newBet model.BetCrash) string {
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
	err = r.db.Model(&model.BetCrash{}).Create(&newBet).Error
	if err != nil {
		return "Хуйня какая-то"
	}
	return "OK"
}

func (r *CrashPostgres) NewCashoutCrash(gameID int, userID string, userMultiplier float64) string {
	err := r.db.Model(&model.BetCrash{}).Where("game_id = ? AND user_id = ?", gameID, userID).Update("user_multiplier", userMultiplier).Error
	if err != nil {
		return "Pizda!"
	}
	return "OK"
}

func (r *CrashPostgres) UpdateWinMultipliers(gameID int, winMultiplier float64) string {
	err := r.db.Model(&model.BetCrash{}).Where("game_id = ?", gameID).Update("win_multiplier", winMultiplier).Error
	if err != nil {
		return "Pizda!"
	}
	return "OK"
}

func (r *CrashPostgres) CreditingWinnings(gameID int) string {
	var bets []model.BetCrash
	if err := r.db.Model(&model.BetCrash{}).Where("game_id = ?", gameID).Find(&bets).Error; err != nil {
		return "Pizda!"
	}
	for _, bet := range bets {
		if bet.UserMultiplier <= bet.WinMultiplier {
			winAmount := bet.Amount * bet.UserMultiplier
			if err := r.db.Model(&model.User{}).Where("id = ?", bet.UserID).Update("balance", gorm.Expr("balance + ?", winAmount)); err != nil {
				return "Pizda!"
			}
		}
	}
	return "OK"
}
