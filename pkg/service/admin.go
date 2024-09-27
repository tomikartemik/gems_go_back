package service

import (
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"github.com/gofrs/uuid"
)

type AdminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) CreateAdmin(admin model.Admin) error {
	admin.Id = uuid.Must(uuid.NewV4()).String()
	admin.Password = GeneratePasswordHash(admin.Password)

	return s.repo.CreateAdmin(admin)
}

func (s *AdminService) SignInAdmin(email string, password string) (string, error) {
	admin, err := s.repo.SignInAdmin(email, GeneratePasswordHash(password))

	if err != nil {
		return "", err
	}

	token := CreateToken(admin.Id)

	return token, nil
}
