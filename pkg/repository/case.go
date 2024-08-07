package repository

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sort"
)

type CasePostgres struct {
	db *gorm.DB
}

func NewCasePostgres(db *gorm.DB) *CasePostgres {
	return &CasePostgres{db: db}
}

func (r *CasePostgres) CreateCase(newCase schema.CaseInput) (int, error) {
	var caseInfo model.Case
	caseInfo.Name = newCase.Name
	caseInfo.Price = newCase.Price
	caseInfo.PhotoLink = newCase.PhotoLink
	caseInfo.Color = newCase.Color
	result := r.db.Model(&model.Case{}).Create(&caseInfo)
	if result.Error != nil {
		return 0, result.Error
	}
	id := caseInfo.ID
	return id, nil
}

func (r *CasePostgres) AddItemsToCase(id int, caseItems []model.CaseItem) (model.Case, error) {
	var caseInfo model.Case
	for i := range caseItems {
		caseItems[i].CaseID = id
	}
	result := r.db.Model(&model.CaseItem{}).Create(&caseItems)
	if result.Error != nil {
		return caseInfo, result.Error
	}
	return caseInfo, nil
}

func (r *CasePostgres) DeleteItemsFromCase(id int) error {
	result := r.db.Where("case_id = ?", id).Delete(&model.CaseItem{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CasePostgres) GetCaseInfo(id int) (model.Case, error) {
	var caseInfo model.Case
	result := r.db.Model(&model.Case{}).Where("id = ?", id).First(&caseInfo)
	if result.Error != nil {
		return caseInfo, result.Error
	}
	return caseInfo, nil
}

func (r *CasePostgres) GetCaseItems(caseId int) ([]model.ItemWithID, error) {
	var caseItems []model.CaseItem
	if err := r.db.Where("case_id = ?", caseId).Find(&caseItems).Error; err != nil {
		return nil, err
	}
	var itemsWithID []model.ItemWithID
	var item model.ItemWithID
	for i := range caseItems {
		itemId := caseItems[i].ItemID
		if err := r.db.Model(model.Item{}).Where("id = ?", itemId).First(&item).Error; err != nil {
			return nil, err
		}
		itemsWithID = append(itemsWithID, item)
	}
	sort.Slice(itemsWithID, func(i, j int) bool {
		return itemsWithID[i].Rarity < itemsWithID[j].Rarity
	})
	return itemsWithID, nil
}

func (r *CasePostgres) GetAllCases() ([]schema.CaseInfo, error) {
	var caseList []schema.CaseInfo
	if err := r.db.Model(model.Case{}).Find(&caseList).Error; err != nil {
		return nil, err
	}
	return caseList, nil
}

func (r *CasePostgres) UpdateCase(id int, newCase schema.CaseInput) (schema.ShowCase, error) {
	var updatedCase schema.ShowCase
	if err := r.db.Model(&model.Case{}).Where("id = ?", id).Updates(newCase).Error; err != nil {
		return updatedCase, err
	}
	return updatedCase, nil
}

func (r *CasePostgres) DeleteCase(id int) error {
	result := r.db.Where("id = ?", id).Delete(&model.Case{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CasePostgres) CheckThePossibilityOfPurchasing(userId string, caseId int) bool {
	var user model.User
	var caseInfo model.Case

	// Получение информации о пользователе
	if err := r.db.Model(&model.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		return false
	}

	// Получение информации о кейсе
	if err := r.db.Model(&model.Case{}).Where("id = ?", caseId).First(&caseInfo).Error; err != nil {
		return false
	}

	// Проверка достаточности баланса
	if float64(caseInfo.Price) <= user.Balance && user.Balance > 0 {
		return true
	}

	return false
}

func (r *CasePostgres) GetChosenItem(id int) (model.ItemWithID, error) {
	var chosenItem model.ItemWithID
	result := r.db.Model(&model.Item{}).Where("id = ?", id).First(&chosenItem)
	if result.Error != nil {
		return chosenItem, result.Error
	}
	return chosenItem, nil
}

func (r *CasePostgres) AddItemToInventoryAndChangeBalance(userId string, itemId int, caseId int) (int, error) {
	var userItem model.UserItem
	var purchasedCase model.Case
	var newItem, bestItem model.Item
	var user model.User

	// Начало транзакции
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Добавление предмета в инвентарь пользователя
	userItem.ItemID = itemId
	userItem.UserID = userId
	userItem.Sold = false
	if err := tx.Create(&userItem).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Получение информации о пользователе с блокировкой записи
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Получение информации о кейсе
	if err := tx.Model(&model.Case{}).Where("id = ?", caseId).First(&purchasedCase).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Проверка и обновление лучшего предмета пользователя
	if user.BestItemId != 0 {
		if err := tx.Model(&model.Item{}).Where("id = ?", itemId).First(&newItem).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		if err := tx.Model(&model.Item{}).Where("id = ?", user.BestItemId).First(&bestItem).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		if newItem.Price > bestItem.Price {
			if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("best_item_id", newItem.ID).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	} else {
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("best_item_id", itemId).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Обновление баланса пользователя
	if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("balance", gorm.Expr("balance - ?", purchasedCase.Price)).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Коммит транзакции
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return userItem.ID, nil
}
