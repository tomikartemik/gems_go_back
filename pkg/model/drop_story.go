package model

type DropRecord struct {
	ID     int  `json:"id" db:"id" gorm:"autoIncrement"`
	ItemID int  `json:"item_id"`
	Dirty  bool `json:"dirty"`
}
