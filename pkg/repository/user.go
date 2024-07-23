package repository

import (
	"errors"
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
	err := r.db.Where("email = ?", mail).First(&user).Error
	if err != nil {
		return showUser, errors.New("Пользователя с такой почтой не существует!")
	}
	err = r.db.Where("email = ? AND password = ?", mail, password).First(&user).Error
	if err != nil {
		return showUser, errors.New("Неверный пароль!")
	}
	showUser, err = r.GetUserById(user.Id)
	if err != nil {
		return showUser, err
	}
	return showUser, nil
}

func (r *UserPostgres) GetUserById(id string) (schema.ShowUser, error) {
	var user model.User
	var userResponse schema.ShowUser
	err := r.db.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return userResponse, err
	}
	userResponse = schema.ShowUser{
		ID:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		IsActive: user.IsActive,
		Balance:  user.Balance,
		BestItem: model.ItemWithID{},
		Photo:    user.Photo,
	}
	var bestItem model.ItemWithID
	if user.BestItemId != 0 {
		err = r.db.Model(&model.Item{}).Where("id = ?", user.BestItemId).First(&bestItem).Error
	}
	userResponse.BestItem = bestItem
	return userResponse, nil
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

func (r *UserPostgres) GetUserInventory(userId string) ([]model.ItemOfInventory, error) {
	var response []model.ItemOfInventory
	var userItems []model.UserItem

	err := r.db.Model(model.UserItem{}).Where("user_id = ?", userId).Order("id").Find(&userItems).Error
	if err != nil {
		return response, err
	}

	var item model.ItemWithID
	for i := range userItems {
		itemId := userItems[i].ItemID
		if err := r.db.Model(model.Item{}).Where("id = ?", itemId).First(&item).Error; err != nil {
			return nil, err
		}
		response = append(response, model.ItemOfInventory{
			item.ID,
			item.Name,
			item.Rarity,
			item.Price,
			item.PhotoLink,
			item.Color,
			userItems[i].ID,
		})
	}

	return response, nil
}

func (r *UserPostgres) SellItem(userId string, user_item_id int) error {
	var user_item model.UserItem
	var item model.Item

	if err := r.db.Model(&model.UserItem{}).Where("id = ?", user_item_id).Find(&user_item).Error; err != nil {
		return err
	}
	if err := r.db.Model(&model.Item{}).Where("id = ?", user_item.ItemID).Find(&item).Error; err != nil {
		return err
	}
	if err := r.db.Model(&model.User{}).Where("id = ?", userId).Update("balance", gorm.Expr("balance + ?", item.Price)).Error; err != nil {
		return err
	}
	if err := r.db.Model(&model.UserItem{}).Where("id = ?", user_item_id).Delete(user_item).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserPostgres) SellAllItem(userId string) error {
	var userItems []model.UserItem
	var item model.Item
	totalPrice := 0
	err := r.db.Model(&model.UserItem{}).Where("user_id = ?", userId).Find(&userItems).Error
	if err != nil {
		return err
	}
	for _, userItem := range userItems {
		err = r.db.Model(&model.Item{}).Where("id = ?", userItem.ItemID).First(&item).Error
		if err != nil {
			return err
		}
		totalPrice += item.Price
		item = model.Item{}
	}
	err = r.db.Model(&model.UserItem{}).Delete(userItems).Error
	if err != nil {
		return err
	}
	err = r.db.Model(&model.User{}).Where("id = ?", userId).Update("balance", gorm.Expr("balance + ?", totalPrice)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPostgres) ChangeAvatar(userId string, newPhoto int) error {
	err := r.db.Model(&model.User{}).Where("id = ?", userId).Update("balance", gorm.Expr("balance - ?", 10)).Error
	if err != nil {
		return err
	}
	r.db.Model(&model.User{}).Where("id = ?", userId).Update("photo", newPhoto)
	return nil
}
