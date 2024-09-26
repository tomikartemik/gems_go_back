package repository

import (
	"fmt"
	"gems_go_back/pkg/model"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Item{},
		&model.Case{},
		&model.CaseItem{},
		&model.UserItem{},
		&model.BetCrash{},
		&model.CrashRecord{},
		&model.RouletteRecord{},
		&model.BetRoulette{},
		&model.Replenishment{},
		&model.Withdraw{},
		&model.Online{},
		&model.DropRecord{},
		&model.Price{},
		&model.Promo{},
		&model.PromoUsage{},
		&model.Admin{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
