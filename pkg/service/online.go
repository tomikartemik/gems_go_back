package service

import (
	"gems_go_back/pkg/repository"
	"math/rand"
	"time"
)

type OnlineService struct {
	repo repository.Online
}

func NewOnlineService(repo repository.Online) *OnlineService {
	return &OnlineService{repo: repo}
}

func (s *OnlineService) GetOnline() int {
	return s.repo.GetOnline()
}

func (s *OnlineService) SetOnline() {
	rand.Seed(time.Now().UnixNano())
	onlineUsers := rand.Intn(201) - 100 + s.repo.GetOnline()
	s.repo.SetOnline(onlineUsers)
}
