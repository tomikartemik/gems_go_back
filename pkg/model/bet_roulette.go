package model

type BetRoulette struct {
	RouletteID int     `json:"-" db:"roulette_id" gorm:"autoIncrement" gorm:"autoIncrement"`
	GameId     int     `json:"game_id"`
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	UserCell   int     `json:"user_cell"`
	WinCell    int     `json:"win_cell"`
}
