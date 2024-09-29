package service

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
)

type CardService struct {
	repo repository.Card
}

func NewCardService(repo repository.Card) *CardService {
	return &CardService{repo: repo}
}

func (s *CardService) UpdateCard(card model.Card) error {
	card.Type = "main"
	return s.repo.UpdateCard(card)
}

func (s *CardService) GetCard() (model.Card, error) {
	return s.repo.GetCard()
}
