package model

type Promo struct {
	ID     int     `json:"-" db:"id" gorm:"autoIncrement"`
	Promo  string  `json:"promo" db:"promo"`
	Reward float64 `json:"reward" db:"reward"`
}

type PromoUsage struct {
	ID      int    `json:"-" db:"id" gorm:"autoIncrement"`
	PromoID int    `json:"promo_id" db:"promo_id"`
	UserID  string `json:"user_id" db:"user_id"`
}
