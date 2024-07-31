package model

type Case struct {
	ID        int        `json:"-" db:"id" gorm:"autoIncrement"`
	Name      string     `json:"name"`
	Price     int        `json:"price"`
	Items     []CaseItem `gorm:"foreignKey:CaseID"`
	PhotoLink string     `json:"photo_link"`
	Color     string     `json:"color"`
}
