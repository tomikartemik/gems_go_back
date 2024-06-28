package model

type CrashRecord struct {
	ID            int     `json:"id" db:"id" gorm:"autoIncrement"`
	WinMultiplier float64 `json:"win_multiplier"`
}
