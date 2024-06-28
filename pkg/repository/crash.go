package repository

import (
	"gems_go_back/pkg/model"
	"gorm.io/gorm"
)

type CrashPostgres struct {
	db *gorm.DB
}

func NewCrashPostgres(db *gorm.DB) *CrashPostgres {
	return &CrashPostgres{db: db}
}

func (r *CrashPostgres) NewRecord(winMuliplier float64) error {
	record := model.CrashRecord{
		WinMultiplier: winMuliplier,
	}
	result := r.db.Model(&model.CrashRecord{}).Create(&record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CrashPostgres) GetAllRecords() ([]model.CrashRecord, error) {
	var all_records []model.CrashRecord
	err := r.db.Model(&model.CrashRecord{}).Find(&all_records).Error
	if err != nil {
		return nil, err
	}
	return all_records, nil
}
