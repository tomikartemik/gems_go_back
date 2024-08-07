package model

type Item struct {
	ID             int    `json:"-" db:"id" gorm:"autoIncrement"`
	Name           string `json:"name"`
	Rarity         int    `json:"rarity"`
	Price          int    `json:"price"`
	PhotoLink      string `json:"photo_link"`
	Color          string `json:"color"`
	CanBeWithdrawn bool   `json:"can_be_withdrawn"`
}

type ItemWithID struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Rarity         int    `json:"rarity"`
	Price          int    `json:"price"`
	PhotoLink      string `json:"photo_link"`
	Color          string `json:"color"`
	CanBeWithdrawn bool   `json:"can_be_withdrawn"`
}

type ItemOfInventory struct {
	ItemID         int    `json:"id"`
	Name           string `json:"name"`
	Rarity         int    `json:"rarity"`
	Price          int    `json:"price"`
	PhotoLink      string `json:"photo_link"`
	Color          string `json:"color"`
	UserItemID     int    `json:"user_item_id"`
	Sold           bool   `json:"sold"`
	CanBeWithdrawn bool   `json:"can_be_withdrawn"`
}
