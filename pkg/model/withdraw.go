package model

import (
	"time"
)

type Withdraw struct {
	ID           int       `json:"-" db:"id"  gorm:"autoIncrement"`
	UserId       string    `json:"user_id"`
	Username     string    `json:"username"`
	AccountEmail string    `json:"account_email"`
	ItemId       int       `json:"item_id"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
	Code         int       `json:"code"`
}
