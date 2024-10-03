package model

type FakeBets struct {
	ID    int    `json:"id" gorm:"autoIncrement"`
	Name  string `json:"name"`
	Photo int    `json:"photo"`
}
