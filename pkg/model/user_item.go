package model

type UserItem struct {
	ID     int    `json:"-" db:"id" gorm:"primaryKey"`
	UserID string `json:"user_id"`
	ItemID int    `json:"item_id"`
	Item   Item   `json:"item" gorm:"foreignKey:ItemID"`
}
