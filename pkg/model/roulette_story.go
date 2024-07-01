package model

type RouletteRecord struct {
	ID      int `json:"id" db:"id" gorm:"autoIncrement"`
	WinCell int `json:"win_cell"`
}
