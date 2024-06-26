package repository

import "gorm.io/gorm"

type RoulettePostgres struct {
	db *gorm.DB
}

func NewRoulettePostgres(db *gorm.DB) *RoulettePostgres {
	return &RoulettePostgres{db: db}
}

func (r *RoulettePostgres) EchoMSGRoulette(msg string) string {
	return msg
}
