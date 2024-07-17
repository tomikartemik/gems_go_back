package model

import (
	"time"
)

type Withdraw struct {
	ID           int       `json:"-" db:"id"  gorm:"autoIncrement"`
	UserId       string    `json:"user_id"`
	Username     string    `json:"username"`
	AccountEmail string    `json:"account_email"`
	Amount       int       `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
	Completed    bool      `json:"completed"`
	CompletedAt  time.Time `json:"completed_at"`
	Code         int       `json:"code"`
}
