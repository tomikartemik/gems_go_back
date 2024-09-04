package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gems_go_back/pkg/repository"
	"os"
	"strconv"
)

type Endpoint struct {
	URL string `json:"url"`
}

type ReportTo struct {
	Endpoints []Endpoint `json:"endpoints"`
	Group     string     `json:"group"`
	MaxAge    int        `json:"max_age"`
}

type ReplenishmentService struct {
	repo repository.Replenishment
}

func NewReplenishmentService(repo repository.Replenishment) *ReplenishmentService {
	return &ReplenishmentService{repo: repo}
}

func hashData(data string, apiKey string) string {
	hash := hmac.New(sha256.New, []byte(apiKey))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// Функция для создания подписи
func createSignature(secret1, amount, currency, email, i, ip, nonce, shopId string) string {
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", amount, currency, email, i, ip, nonce, shopId)
	fmt.Println(data)
	return hashData(data, secret1)
}

// Функция для отправки запроса на пополнение
func createPaymentRequest(shopID, secret1, secret2, amount, currency, orderID, paymentMethod, email string) (string, error) {

	signature := createSignature(secret1, amount, currency, email, paymentMethod, "20.21.27.109", orderID, shopID)
	url := fmt.Sprintf("https://pay.freekassa.com?currency=%s&email=%s&i=%s&shopId=%s&nonce=%s&amount=%s&signature=%s&ip=20.21.27.109", currency, email, paymentMethod, shopID, orderID, amount, signature)
	fmt.Println(url)
	return url, nil
}

func (s *ReplenishmentService) NewReplenishment(userId string, amount float64, promo string) (string, error) {
	var replenishmentID string
	var email string
	var err error
	rewardInfo, err := s.repo.GetReward(promo, userId)
	if err != nil {
		replenishmentID, email, err = s.repo.NewReplenishment(userId, amount)
	} else {
		replenishmentID, email, err = s.repo.NewReplenishment(userId, amount*rewardInfo)
	}

	var merchantID = os.Getenv("MERCHANT_ID")
	var secret1 = os.Getenv("SECRET_1")
	var secret2 = os.Getenv("SECRET_2")
	var currency = os.Getenv("CURRENCY")

	location, err := createPaymentRequest(merchantID, secret1, secret2, strconv.FormatFloat(amount, 'f', -1, 64), currency, replenishmentID, "36", email)
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
