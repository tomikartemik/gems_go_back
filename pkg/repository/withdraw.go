package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
	"time"
)

type WithdrawPostgres struct {
	db *gorm.DB
}

func NewWithdrawPostgres(db *gorm.DB) *WithdrawPostgres {
	return &WithdrawPostgres{db: db}
}

func (r *WithdrawPostgres) CreateWithdraw(withdraw model.Withdraw) (model.Withdraw, error) {
	var user model.User
	var newWithdraw model.Withdraw
	err := r.db.Model(&model.User{}).Where("id = ?", withdraw.UserId).First(&user).Error
	if err != nil {
		return newWithdraw, err
	}
	err = r.db.Model(&model.User{}).Where("id = ?", withdraw.UserId).Update("balance", gorm.Expr("balance - ?", float64(withdraw.Amount))).Error
	if err != nil {
		return newWithdraw, err
	}
	newWithdraw.UserId = withdraw.UserId
	newWithdraw.Username = user.Username
	newWithdraw.AccountEmail = withdraw.AccountEmail
	newWithdraw.Code = withdraw.Code
	newWithdraw.Amount = withdraw.Amount
	newWithdraw.CreatedAt = time.Now()
	err = r.db.Model(&model.Withdraw{}).Create(&newWithdraw).Error
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
	err := r.db.Model(&model.Withdraw{}).Where("id = ?", withdrawId).Update("completed", true).Error
	if err != nil {
		return err
	}
	err = r.db.Model(&model.Withdraw{}).Where("id = ?", withdrawId).Update("completed_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}
