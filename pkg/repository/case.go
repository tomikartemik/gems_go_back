package repository

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/schema"
	"gorm.io/gorm"
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

func (r *CasePostgres) GetItemsWithWeights(id int) ([]model.CaseItem, error) {
	var caseItems []model.CaseItem
	if err := r.db.Where("case_id = ?", id).Find(&caseItems).Error; err != nil {
		return nil, err
	}
	return caseItems, nil
}

func (r *CasePostgres) GetChosenItem(id int) (model.ItemWithID, error) {
	var chosenItem model.ItemWithID
	result := r.db.Model(&model.Item{}).Where("id = ?", id).First(&chosenItem)
	if result.Error != nil {
		return chosenItem, result.Error
	}
	return chosenItem, nil
}

func (r *CasePostgres) NewCaseRecord(case_id int) error {
	record := model.CaseRecord{CaseID: case_id}
	if err := r.db.Model(&model.CaseRecord{}).Create(record).Error; err != nil {
		return err
	}
	return nil
}

func (r *CasePostgres) GetAllCaseRecords() ([]schema.CaseInfo, error) {
	var lastCases []schema.CaseInfo
	var caseRecords []model.CaseRecord
	var caseInfo schema.CaseInfo
	err := r.db.Model(&model.CaseRecord{}).Order("id desc").Limit(10).Find(&caseRecords).Error
	if err != nil {
		return lastCases, err
	}
	for _, currCase := range caseRecords {
		err = r.db.Model(&model.Case{}).Where("id = ?", currCase.CaseID).Find(caseInfo).Error
		if err != nil {
			return lastCases, err
		}
		lastCases = append(lastCases, caseInfo)
	}
	return lastCases, nil
}
