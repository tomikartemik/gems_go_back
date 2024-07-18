package model

type CaseItem struct {
	ID     int  `json:"-" db:"id" gorm:"primaryKey"`
	CaseID int  `json:"case_id"`
	ItemID int  `json:"item_id"`
	Item   Item `json:"item" gorm:"foreignKey:ItemID"`
}
