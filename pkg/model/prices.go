package model

type Price struct {
	ID       int     `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Position string  `json:"position"`
	Price    float64 `json:"price"`
}
