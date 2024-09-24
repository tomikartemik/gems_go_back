package service

import (
	"fmt"
	"gems_go_back/pkg/repository"
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
)

type ReceiptService struct {
	repo repository.Receipt
}

func NewReceiptService(repo repository.Receipt) *ReceiptService {
	return &ReceiptService{repo: repo}
}

func (s *ReceiptService) UploadReceipt(file *multipart.FileHeader) (string, error) {
	fileExt := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + fileExt
	pathToFile := filepath.Join("./uploads", newFileName)

	err := s.repo.UploadReceipt(file, pathToFile)
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFileName)
	return fileURL, nil
}
