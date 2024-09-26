package model

type Admin struct {
	Id       string `json:"-" db:"id" gorm:"primaryKey"`
	Username string `json:"username" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required"`
}
