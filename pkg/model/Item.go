package model

type Item struct {
	ID     int    `json:"-" db:"id" gorm:"autoIncrement"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
	Price  int    `json:"price"`
}

type ItemWithID struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
	Price  int    `json:"price"`
}
