package service

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
)

type FakeBetService struct {
	repo repository.FakeBets
}

func NewFakeBetService(repo repository.FakeBets) *FakeBetService {
	return &FakeBetService{repo: repo}
}

func (s *FakeBetService) CreateFakeBetter(user model.FakeBets) {
	s.repo.CreateFakeBet(user)
}
