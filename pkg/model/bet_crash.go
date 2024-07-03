package model

type BetCrash struct {
	BetID          int     `json:"-" db:"bet_id"  gorm:"autoIncrement"`
	GameId         int     `json:"game_id"`
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	UserMultiplier float64 `json:"user_multiplier"`
	WinMultiplier  float64 `json:"win_multiplier"`
}
