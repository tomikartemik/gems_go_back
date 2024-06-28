package repository

import (
	"errors"
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/schema"
	"gorm.io/gorm"
)

type UserPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(user model.User) (schema.ShowUser, error) {
	var createdUser schema.ShowUser
	id := user.Id
	result := r.db.Model(&model.User{}).Create(&user)
	if result.Error != nil {
		return createdUser, result.Error
	}
	createdUser, err := r.GetUserById(id)
	if err != nil {
		return createdUser, err
	}
	return createdUser, nil
}

func (r *UserPostgres) SignIn(mail, password string) (schema.ShowUser, error) {
	var user model.User
	var showUser schema.ShowUser
	result := r.db.Where("email = ? AND password = ?", mail, password).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return showUser, errors.New("user not found")
	}
	if result.Error != nil {
		return showUser, result.Error
	}

	showUser, err := r.GetUserById(user.Id)
	if err != nil {
		return showUser, err
	}
	return showUser, nil
}

func (r *UserPostgres) GetUserById(id string) (schema.ShowUser, error) {
	var user schema.ShowUser
	err := r.db.Model(model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func (r *UserPostgres) UpdateUser(id string, user schema.InputUser) (schema.ShowUser, error) {
	if err := r.db.Model(&model.User{}).Where("id = ?", id).Updates(&user).Error; err != nil {
		return schema.ShowUser{}, err
	}

	var updatedUser model.User
	var showUser schema.ShowUser
	result := r.db.Where("id = ?", id).First(&updatedUser)
	if result.Error != nil {
		return showUser, result.Error
	}

	showUser, err := r.GetUserById(id)
	if err != nil {
		return showUser, err
	}
	return showUser, nil
}

func (r *UserPostgres) GetUserInventory(userId string) ([]model.ItemWithID, error) {
	var itemsWithID []model.ItemWithID

	var userItems []model.UserItem
	err := r.db.Model(model.UserItem{}).Where("user_id = ?", userId).Find(&userItems).Error
	fmt.Println(userItems)
	if err != nil {
		return itemsWithID, err
	}

	var item model.ItemWithID
	for i := range userItems {
		itemId := userItems[i].ItemID
		if err := r.db.Model(model.Item{}).Where("id = ?", itemId).First(&item).Error; err != nil {
			return nil, err
		}
		itemsWithID = append(itemsWithID, item)
	}
	return itemsWithID, nil
}

func (r *UserPostgres) AddItemToInventory(userId string, itemId int) (model.UserItem, error) {
	var userItem model.UserItem
	userItem.ItemID = itemId
	userItem.UserID = userId
	result := r.db.Create(&userItem)
	if result.Error != nil {
		return userItem, result.Error
	}
	return userItem, nil
}