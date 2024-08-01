package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
)

type CrashPostgres struct {
	db *gorm.DB
}

func NewCrashPostgres(db *gorm.DB) *CrashPostgres {
	return &CrashPostgres{db: db}
}

func (r *CrashPostgres) NewCrashRecord(winMuliplier float64) error {
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

	// Начало транзакции
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Получение пользователя с блокировкой записи
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.User{}).Where("id = ?", newBet.UserID).First(&user).Error
	if err != nil {
		tx.Rollback()
		return "User not found!"
	}

	// Проверка баланса
	if user.Balance < newBet.Amount || newBet.Amount == 0 {
		tx.Rollback()
		return "Not enough money!"
	}

	// Обновление баланса
	err = tx.Model(&model.User{}).Where("id = ?", user.Id).Update("balance", gorm.Expr("balance - ?", newBet.Amount)).Error
	if err != nil {
		tx.Rollback()
		return "Error updating balance"
	}

	// Создание записи о ставке
	err = tx.Model(&model.BetCrash{}).Create(&newBet).Error
	if err != nil {
		tx.Rollback()
		return "Error creating bet"
	}

	// Коммит транзакции
	if err := tx.Commit().Error; err != nil {
		return "Transaction commit error"
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

func (r *CrashPostgres) CreditingWinningsCrash(gameID int) string {
	var bets []model.BetCrash
	if err := r.db.Model(&model.BetCrash{}).Where("game_id = ?", gameID).Find(&bets).Error; err != nil {
		return "Pizda!"
	}
	for _, bet := range bets {
		if bet.UserMultiplier <= bet.WinMultiplier {
			winAmount := math.Round(bet.Amount*bet.UserMultiplier*100) / 100
			if err := r.db.Model(&model.User{}).Where("id = ?", bet.UserID).Update("balance", gorm.Expr("balance + ?", winAmount)); err != nil {
				return "Pizda!"
			}
		}
	}
	return "OK"
}

func (r *CrashPostgres) GetUsersPhotoAndNickForCrash(userId string) (string, error) {
	var user model.User
	if err := r.db.Model(&model.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		return "", err
	}
	return user.Username, nil
}
