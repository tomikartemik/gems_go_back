package service

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"github.com/google/uuid"
	"math"
	"mime/multipart"
	"path/filepath"
)

type OwnReplenishmentService struct {
	repo        repository.OwnReplenishment
	repoReceipt repository.Receipt
	repoUser    repository.User
}

func NewOwnReplenishmentService(repo repository.OwnReplenishment, repoReceipt repository.Receipt, repoUser repository.User) *OwnReplenishmentService {
	return &OwnReplenishmentService{repo: repo, repoReceipt: repoReceipt, repoUser: repoUser}
}

func (s *OwnReplenishmentService) CreateReplenishment(amount float64, userID string, file *multipart.FileHeader) error {
	replenishment := model.OwnReplenishment{}

	receiptURL, err := s.uploadReceipt(file)

	if err != nil {
		return err
	}

	replenishment.ReceiptURL = receiptURL
	replenishment.Amount = amount
	replenishment.UserId = userID
	replenishment.Status = "Processing"
	err = s.repo.CreateReplenishment(replenishment)
	if err != nil {
		return err
	}

	return nil
}

func (s *OwnReplenishmentService) uploadReceipt(file *multipart.FileHeader) (string, error) {
	fileExt := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + fileExt
	pathToFile := filepath.Join("./uploads", newFileName)

	err := s.repoReceipt.UploadReceipt(file, pathToFile)
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFileName)
	return fileURL, nil
}

func (s *OwnReplenishmentService) GetReplenishments(sortOrder string, page int) (model.OwnReplenishmentOutput, error) {
	ownReplOutput := model.OwnReplenishmentOutput{}
	lastID, err := s.repo.GetLastId()
	if err != nil {
		return model.OwnReplenishmentOutput{}, err
	}
	pagesCount := int(math.Ceil(float64(lastID) / float64(10)))
	repls, err := s.repo.GetReplenishments(sortOrder, page)
	if err != nil {
		return ownReplOutput, err
	}
	var responses []model.OwnReplenishmentsResponse

	for _, repl := range repls {
		user, err := s.repoUser.GetUserById(repl.UserId)
		username := user.Username
		if err != nil {
			username = ""
		}

		response := model.OwnReplenishmentsResponse{
			ID:         repl.ID,
			UserId:     repl.UserId,
			Username:   username,
			Amount:     repl.Amount,
			ReceiptURL: repl.ReceiptURL,
			Status:     repl.Status,
		}

		responses = append(responses, response)
	}

	ownReplOutput = model.OwnReplenishmentOutput{
		Replenishments: responses,
		PagesCount:     pagesCount,
	}

	return ownReplOutput, nil
}

func (s *OwnReplenishmentService) ChangeStatus(replenishmentID int, status string) error {
	replenishment, err := s.repo.ChangeStatus(replenishmentID, status)
	if err != nil {
		return err
	}
	if status == "Approved" {
		err = s.repo.ChangeBalance(replenishment.UserId, replenishment.Amount)
		if err != nil {
			return err
		}
	}
	return nil
}
