package repository

import "gorm.io/gorm"

type CrashPostgres struct {
	db *gorm.DB
}

func NewCrashPostgres(db *gorm.DB) *CrashPostgres {
	return &CrashPostgres{db: db}
}

func (r *CrashPostgres) EchoMSG(msg string) string {
	return msg
}
