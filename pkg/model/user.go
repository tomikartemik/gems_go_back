package model

type User struct {
	Id         string     `json:"-" db:"id" gorm:"primaryKey"`
	Email      string     `json:"email" binding:"required" gorm:"unique"`
	Username   string     `json:"username" binding:"required" gorm:"unique"`
	Password   string     `json:"password" binding:"required"`
	Balance    float64    `json:"balance"`
	IsActive   bool       `json:"is_active"`
	IsAdmin    bool       `json:"is_admin"`
	BestItemId int        `json:"best_item_id"`
	Items      []UserItem `json:"items" gorm:"foreignKey:UserID"`
}
