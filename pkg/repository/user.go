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
	var user model.User
	var userResponse schema.ShowUser
	err := r.db.Model(model.User{}).Where("id = ?", id).First(&user).Error
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
	}
	var bestItem model.ItemWithID
	err = r.db.Model(&model.Item{}).Where("id = ?", user.BestItemId).First(&bestItem).Error
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
	err := r.db.Model(model.UserItem{}).Where("user_id = ?", userId).Find(&userItems).Error
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
			userItems[i].ID,
		})
	}

	return response, nil
}

//func (r *UserPostgres) AddItemToInventory(userId string, itemId int) (model.UserItem, error) {
//	var userItem model.UserItem
//	userItem.ItemID = itemId
//	userItem.UserID = userId
//	result := r.db.Create(&userItem)
//	if result.Error != nil {
//		return userItem, result.Error
//	}
//	var newItem, bestItem model.Item
//	var user model.User
//	if err := r.db.Model(&model.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
//		return userItem, nil
//	}
//	if user.BestItemId != 0 {
//		if err := r.db.Model(&model.Item{}).Where("id = ?", itemId).First(&newItem).Error; err != nil {
//			return userItem, nil
//		}
//		if err := r.db.Model(&model.Item{}).Where("id = ?", user.BestItemId).First(&bestItem).Error; err != nil {
//			return userItem, nil
//		}
//		if newItem.Price > bestItem.Price {
//			if err := r.db.Model(&model.User{}).Where("id = ?", userId).Update("best_item_id", newItem.ID).Error; err != nil {
//				return userItem, nil
//			}
//		}
//	} else {
//		if err := r.db.Model(&model.User{}).Where("id = ?", userId).Update("best_item_id", itemId).Error; err != nil {
//			return userItem, nil
//		}
//	}
//	return userItem, nil
//}

func (r *UserPostgres) SellAnItem(userId string, user_item_id int) error {
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
