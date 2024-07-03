package model

type Replenishment struct {
	ID      int     `json:"-" db:"id"  gorm:"autoIncrement"`
	UserID  string  `json:"user_id"`
	Amount  float64 `json:"amount"`
	Success bool    `json:"success"`
}
