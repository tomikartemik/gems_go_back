package service

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"gems_go_back/pkg/schema"
	"math/rand"
	"time"
)

type CaseService struct {
	repo repository.Case
}

func NewCaseService(repo repository.Case) *CaseService {
	return &CaseService{repo: repo}
}

func (s *CaseService) CreateCase(newCase model.Case) (schema.ShowCase, error) {
	var createdCase schema.ShowCase
	var caseInput schema.CaseInput
	caseInput.Name = newCase.Name
	caseInput.Price = newCase.Price
	createdCaseId, err := s.repo.CreateCase(caseInput)
	if err != nil {
		return createdCase, err
	}
	var items []model.CaseItem
	items = newCase.Items
	_, err = s.AddItemsToCase(createdCaseId, items)
	if err != nil {
		return createdCase, err
	}
	createdCase, err = s.GetCase(createdCaseId)
	if err != nil {
		return createdCase, err
	}
	return createdCase, nil
}

func (s *CaseService) GetCase(id int) (schema.ShowCase, error) {
	var caseOutput schema.ShowCase
	caseInfo, err := s.repo.GetCaseInfo(id)
	if err != nil {
		return caseOutput, err
	}
	caseItems, err := s.repo.GetCaseItems(id)
	if err != nil {
		return caseOutput, err
	}
	caseOutput.Id = caseInfo.ID
	caseOutput.Name = caseInfo.Name
	caseOutput.Price = caseInfo.Price
	caseOutput.Items = caseItems
	return caseOutput, nil
}

func (s *CaseService) GetAllCases() ([]schema.CaseInfo, error) {
	return s.repo.GetAllCases()
}

func (s *CaseService) AddItemsToCase(id int, caseItems []model.CaseItem) (schema.ShowCase, error) {
	var caseOutput schema.ShowCase
	_, err := s.repo.AddItemsToCase(id, caseItems)
	if err != nil {
		return caseOutput, err
	}
	caseOutput, err = s.GetCase(id)
	if err != nil {
		return caseOutput, err
	}
	return caseOutput, nil
}

func (s *CaseService) GetCaseItems(caseId int) ([]model.ItemWithID, error) {
	return s.repo.GetCaseItems(caseId)
}

func (s *CaseService) UpdateCase(id int, updates model.Case) (schema.ShowCase, error) {
	var infoForUpdate schema.CaseInput
	infoForUpdate.Name = updates.Name
	infoForUpdate.Price = updates.Price

	var updatedItems []model.CaseItem
	updatedItems = updates.Items

	var updatedCaseInfo schema.ShowCase
	updatedCaseInfo, err := s.repo.UpdateCase(id, infoForUpdate)
	if err != nil {
		return updatedCaseInfo, err
	}

	err = s.repo.DeleteItemsFromCase(id)
	if err != nil {
		return updatedCaseInfo, err
	}

	_, err = s.AddItemsToCase(id, updatedItems)
	if err != nil {
		return updatedCaseInfo, err
	}
	updatedCaseInfo, err = s.GetCase(id)
	if err != nil {
		return updatedCaseInfo, err
	}
	return updatedCaseInfo, nil
}

func (s *CaseService) DeleteCase(caseId int) error {
	err := s.repo.DeleteItemsFromCase(caseId)
	if err != nil {
		return err
	}

	err = s.repo.DeleteCase(caseId)
	if err != nil {
		return err
	}

	return nil
}

func (s *CaseService) OpenCase(userId string, caseId int) (model.ItemWithID, error) {
	var chosenItem model.ItemWithID

	var caseItems []model.CaseItem
	caseItems, err := s.repo.GetItemsWithWeights(caseId)
	if err != nil {
		return chosenItem, err
	}

	//Суммируем все веса
	totalWeightInCase := 0
	for _, caseItem := range caseItems {
		totalWeightInCase += caseItem.Weight
	}

	//Генерируем случайное число
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(totalWeightInCase)

	//Выбираем случайно предмет, учитывая веса
	for _, caseItem := range caseItems {
		if randomNumber < caseItem.Weight {
			id := caseItem.ItemID
			chosenItem, err = s.repo.GetChosenItem(id)
			if err != nil {
				return chosenItem, err
			}
			return chosenItem, nil
		}
		randomNumber -= caseItem.Weight
	}
	err = s.repo.NewCaseRecord(caseId)
	if err != nil {
		fmt.Println(err)
		return model.ItemWithID{}, err
	}
	err = s.repo.AddItemToInventoryAndChangeBalance(userId, chosenItem.ID)
	if err != nil {
		fmt.Println(err)
		return model.ItemWithID{}, err
	}
	return chosenItem, nil
}

func (s *CaseService) GetAllCaseRecords() ([]schema.CaseInfo, error) {
	var records []schema.CaseInfo
	records, err := s.repo.GetAllCaseRecords()
	if err != nil {
		return records, err
	}
	return records, nil
}
