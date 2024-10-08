package service

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"gems_go_back/pkg/schema"
	"math/rand"
	"sort"
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
	caseInput.PhotoLink = newCase.PhotoLink
	if newCase.Price <= 100 {
		caseInput.Color = "#ffffff"
	} else if newCase.Price <= 200 {
		caseInput.Color = "#db2f4c"
	} else if newCase.Price <= 400 {
		caseInput.Color = "#00989E"
	} else if newCase.Price <= 500 {
		caseInput.Color = "#6cbf01"
	} else {
		caseInput.Color = "#fd9d2d"
	}
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
	caseOutput.PhotoLink = caseInfo.PhotoLink
	caseOutput.Collection = caseInfo.Collection
	caseOutput.Items = caseItems
	return caseOutput, nil
}

func (s *CaseService) GetAllCases() ([]schema.CaseInfo, error) {
	allCases, err := s.repo.GetAllCases()
	if err != nil {
		return allCases, err
	}

	sort.Slice(allCases, func(i, j int) bool {
		return allCases[i].Id < allCases[j].Id
	})

	return allCases, nil
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

func (s *CaseService) OpenCase(userId string, caseId int) (model.ItemWithID, int, error) {
	if s.repo.CheckThePossibilityOfPurchasing(userId, caseId) == false {
		return model.ItemWithID{}, -1, nil
	}
	var chosenItem model.ItemWithID
	var caseItems []model.ItemWithID
	caseItems, err := s.repo.GetCaseItems(caseId)
	if err != nil {
		return chosenItem, 0, err
	}
	totalWeightInCase := 0
	for _, caseItem := range caseItems {
		totalWeightInCase += caseItem.Rarity
	}

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(totalWeightInCase)
	for _, caseItem := range caseItems {
		if randomNumber < caseItem.Rarity {
			id := caseItem.ID
			chosenItem, err = s.repo.GetChosenItem(id)
			if err != nil {
				break
			}
			break
		}
		randomNumber -= caseItem.Rarity
	}
	if err != nil {
		fmt.Println(err)
		return model.ItemWithID{}, 0, err
	}
	userItemId, err := s.repo.AddItemToInventoryAndChangeBalance(userId, chosenItem.ID, caseId)
	if err != nil {
		fmt.Println(err)
		return model.ItemWithID{}, 0, err
	}
	return chosenItem, userItemId, nil
}
