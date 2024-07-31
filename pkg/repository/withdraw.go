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
	if user.Balance < newWithdraw.Price || user.Balance <= 0 {
		return model.Withdraw{Username: "денег не хватает, броук"}, err
	}
	err = r.db.Model(&model.User{}).Where("id = ?", withdraw.UserId).Update("balance", gorm.Expr("balance - ?", withdraw.Price)).Error
	if err != nil {
		return newWithdraw, err
	}
	newWithdraw.UserId = withdraw.UserId
	newWithdraw.Username = user.Username
	newWithdraw.AccountEmail = withdraw.AccountEmail
	newWithdraw.Code = withdraw.Code
	newWithdraw.Amount = withdraw.Amount
	newWithdraw.Price = withdraw.Price
	newWithdraw.Status = "processing"
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

func (r *WithdrawPostgres) GetPositionPrice(amount int) (float64, error) {
	var postition model.Price
	err := r.db.Model(&model.Price{}).Where("position = ?", amount).Find(&postition).Error
	if err != nil {
		return 0, err
	}
	return postition.Price, nil
}

func (r *WithdrawPostgres) GetPositionPrices() []model.Price {
	var positions []model.Price
	if err := r.db.Model(&model.Price{}).Find(&positions).Error; err != nil {
		return []model.Price{}
	}
	return positions
}
