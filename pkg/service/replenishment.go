package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gems_go_back/pkg/repository"
	"strconv"
)

// Структуры для парсинга JSON
type Endpoint struct {
	URL string `json:"url"`
}

type ReportTo struct {
	Endpoints []Endpoint `json:"endpoints"`
	Group     string     `json:"group"`
	MaxAge    int        `json:"max_age"`
}

var merchantID = "46264"
var secret1 = "@R-m/.IntF(1eh&"
var secret2 = "1YHPU6{azd?M*LE"
var currency = "RUB"

type ReplenishmentService struct {
	repo repository.Replenishment
}

func NewReplenishmentService(repo repository.Replenishment) *ReplenishmentService {
	return &ReplenishmentService{repo: repo}
}

// Функция для генерации MD5 хеша
func generateMD5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Функция для создания подписи
func createSignature(merchantID, secret1, amount, currency, orderID string) string {
	data := fmt.Sprintf("%s:%s:%s:%s:%s", merchantID, amount, secret1, currency, orderID)
	return generateMD5Hash(data)
}

// Функция для отправки запроса на пополнение
func createPaymentRequest(merchantID, secret1, secret2, amount, currency, orderID, description, email string) (string, error) {

	signature := createSignature(merchantID, secret1, amount, currency, orderID)
	url := fmt.Sprintf("https://pay.freekassa.com?currency=%s&email=%s&i=%s&m=%s&o=%s&oa=%s&s=%s", currency, email, description, merchantID, orderID, amount, signature)
	fmt.Println(url)
	return url, nil
}

func (s *ReplenishmentService) NewReplenishment(userId string, amount float64) (string, error) {
	replenishmentID, email, err := s.repo.NewReplenishment(userId, amount)
	if err != nil {
		return "", err
	}
	location, err := createPaymentRequest(merchantID, secret1, secret2, strconv.FormatFloat(amount, 'f', -1, 64), currency, replenishmentID, "test", email)
	if err != nil {
		return "", err
	}
	return location, nil
}

func (s *ReplenishmentService) AcceptReplenishment(replenishmentID int) {
	err := s.repo.AcceptReplenishment(replenishmentID)
	if err != nil {
		fmt.Println(err)
	}
}
