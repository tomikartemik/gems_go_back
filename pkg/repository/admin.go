package repository

import (
	"errors"
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type AdminPostgres struct {
	db *gorm.DB
}

func NewAdminPostgres(db *gorm.DB) *AdminPostgres {
	return &AdminPostgres{db: db}
}

func (r *AdminPostgres) CreateAdmin(admin model.Admin) error {
	err := r.db.Model(&model.Admin{}).Create(&admin).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AdminPostgres) SignInAdmin(mail, password string) (model.Admin, error) {
	var admin model.Admin
	err := r.db.Model(&model.Admin{}).Where("username = ?", mail).First(&admin).Error
	if err != nil {
		return admin, errors.New("Пользователя с такой почтой не существует!")
	}
	err = r.db.Model(&model.Admin{}).Where("username = ? AND password = ?", mail, password).First(&admin).Error
	if err != nil {
		return admin, errors.New("Неверный пароль!")
	}
	return admin, nil
}
