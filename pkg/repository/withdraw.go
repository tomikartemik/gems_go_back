package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type WithdrawPostgres struct {
	db *gorm.DB
}

func NewWithdrawPostgres(db *gorm.DB) *WithdrawPostgres {
	return &WithdrawPostgres{db: db}
}

func (r *WithdrawPostgres) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *WithdrawPostgres) CreateWithdraw(tx *gorm.DB, withdraw model.Withdraw) (model.Withdraw, error) {
	var user model.User
	var newWithdraw model.Withdraw

	// Получение пользователя с блокировкой записи
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.User{}).Where("id = ?", withdraw.UserId).First(&user).Error
	if err != nil {
		return newWithdraw, err
	}

	// Проверка баланса
	if user.Balance < withdraw.Price || user.Balance <= 0 {
		return model.Withdraw{Username: "денег не хватает, броук"}, nil
	}

	// Обновление баланса
	err = tx.Model(&model.User{}).Where("id = ?", withdraw.UserId).Update("balance", gorm.Expr("balance - ?", withdraw.Price)).Error
	if err != nil {
		return newWithdraw, err
	}

	// Создание записи о выводе средств
	newWithdraw = model.Withdraw{
		UserId:       withdraw.UserId,
		Username:     user.Username,
		AccountEmail: withdraw.AccountEmail,
		Code:         withdraw.Code,
		ItemId:       withdraw.ItemId,
		Price:        withdraw.Price,
		Status:       "processing",
		CreatedAt:    time.Now(),
	}
	err = tx.Model(&model.Withdraw{}).Create(&newWithdraw).Error
	if err != nil {
		return newWithdraw, err
	}

	return newWithdraw, nil
}

func (r *WithdrawPostgres) GetWithdraw(withdrawId int) (model.Withdraw, error) {
	var withdraw model.Withdraw
	err := r.db.Model(&model.Withdraw{}).Where("id = ?", withdrawId).First(&withdraw).Error
	if err != nil {
		return model.Withdraw{}, err
	}
	return withdraw, nil
}

func (r *WithdrawPostgres) CompleteWithdraw(withdrawId int) error {
	err := r.db.Model(&model.Withdraw{}).Where("id = ?", withdrawId).Update("status", "completed").Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WithdrawPostgres) CancelWithdraw(withdrawId int) error {
	err := r.db.Model(&model.Withdraw{}).Where("id = ?", withdrawId).Update("status", "canceled").Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WithdrawPostgres) ReturnMoneyBecauseCanceled(currentWithdraw model.Withdraw) {
	r.db.Model(&model.User{}).Where("id = ?", currentWithdraw.UserId).Update("balance", gorm.Expr("balance + ?", currentWithdraw.Price))
}

func (r *WithdrawPostgres) GetUsersWithdraws(userId string) ([]model.Withdraw, error) {
	var withdraws []model.Withdraw
	err := r.db.Model(&model.Withdraw{}).Where("user_id = ?", userId).Find(&withdraws).Error
	if err != nil {
		return withdraws, err
	}
	return withdraws, nil
}

func (r *WithdrawPostgres) GetPositionPrice(itemId int) (int, error) {
	var item model.Item
	err := r.db.Model(&model.Item{}).Where("id = ?", itemId).Find(&item).Error
	if err != nil {
		return 0, err
	}
	return item.Price, nil
}

func (r *WithdrawPostgres) GetPositionPrices() []model.Price {
	var positions []model.Price
	if err := r.db.Model(&model.Price{}).Find(&positions).Error; err != nil {
		return []model.Price{}
	}
	return positions
}
