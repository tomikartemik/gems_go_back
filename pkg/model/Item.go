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

type ItemOfInventory struct {
	ItemID     int    `json:"id"`
	Name       string `json:"name"`
	Rarity     int    `json:"rarity"`
	Price      int    `json:"price"`
	UserItemID int    `json:"user_item_id"`
}
