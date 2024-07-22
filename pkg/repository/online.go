package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type OnlinePostgres struct {
	db *gorm.DB
}

func NewOnlinePostgres(db *gorm.DB) *OnlinePostgres {
	return &OnlinePostgres{db: db}
}

func (r *OnlinePostgres) GetOnline() int {
	var online model.Online
	err := r.db.Model(&model.Online{}).First(&online).Error
	if err != nil {
		return 0
	}
	return online.UsersOnline
}

func (r *OnlinePostgres) SetOnline(usersOnline int) {
	var newOnline model.Online
	newOnline.UsersOnline = usersOnline
	newOnline.Status = "Regular"
	_ = r.db.Model(&model.Online{}).Where("status = ?", "Regular").Updates(&newOnline).Error
}
