package repository

import (
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
)

type ReceiptPostgres struct {
	db *gorm.DB
}

func NewReceiptPostgres(db *gorm.DB) *ReceiptPostgres {
	return &ReceiptPostgres{db: db}
}

func (r *ReceiptPostgres) UploadReceipt(file *multipart.FileHeader, filepath string) error {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err := os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			return err
		}
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
