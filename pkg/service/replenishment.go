package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gems_go_back/pkg/repository"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

//type Endpoint struct {
//	URL string `json:"url"`
//}
//
//type ReportTo struct {
//	Endpoints []Endpoint `json:"endpoints"`
//	Group     string     `json:"group"`
//	MaxAge    int        `json:"max_age"`
//}

type ReplenishmentService struct {
	repo repository.Replenishment
}

func NewReplenishmentService(repo repository.Replenishment) *ReplenishmentService {
	return &ReplenishmentService{repo: repo}
}

//// Функция для генерации MD5 хеша
//func generateMD5Hash(data string) string {
//	hash := md5.Sum([]byte(data))
//	return hex.EncodeToString(hash[:])
//}
//
//// Функция для создания подписи
//func createSignature(merchantID, secret1, amount, currency, orderID string) string {
//	data := fmt.Sprintf("%s:%s:%s:%s:%s", merchantID, amount, secret1, currency, orderID)
//	return generateMD5Hash(data)
//}
//
//// Функция для отправки запроса на пополнение
//func createPaymentRequest(merchantID, secret1, secret2, amount, currency, orderID, description, email string) (string, error) {
//
//	signature := createSignature(merchantID, secret1, amount, currency, orderID)
//	url := fmt.Sprintf("https://pay.freekassa.com?currency=%s&email=%s&i=%s&m=%s&o=%s&oa=%s&s=%s", currency, email, description, merchantID, orderID, amount, signature)
//	fmt.Println(url)
//	return url, nil
//}
//
//func (s *ReplenishmentService) NewReplenishment(userId string, amount float64, promo string) (string, error) {
//	var replenishmentID string
//	var email string
//	var err error
//	rewardInfo, err := s.repo.GetReward(promo, userId)
//	if err != nil {
//		replenishmentID, email, err = s.repo.NewReplenishment(userId, amount)
//	} else {
//		replenishmentID, email, err = s.repo.NewReplenishment(userId, amount*rewardInfo)
//	}
//	var merchantID = os.Getenv("MERCHANT_ID")
//	var secret1 = os.Getenv("SECRET_1")
//	var secret2 = os.Getenv("SECRET_2")
//	var currency = os.Getenv("CURRENCY")
//
//	location, err := createPaymentRequest(merchantID, secret1, secret2, strconv.FormatFloat(amount, 'f', -1, 64), currency, replenishmentID, "12", email)
//	if err != nil {
//		return "", err
//	}
//	return location, nil
//}

func (s *ReplenishmentService) AcceptReplenishment(replenishmentID int) {
	err := s.repo.AcceptReplenishment(replenishmentID)
	if err != nil {
		fmt.Println(err)
	}
}

type CreateOrderRequest struct {
	Amount     float64 `json:"amount"`   // Сумма платежа
	CurrencyID string  `json:"currency"` // Идентификатор валюты
	Email      string  `json:"email"`    // Email покупателя
	ShopID     int     `json:"shopId"`
	I          int     `json:"i"`
	IP         string  `json:"ip"`        // IP адрес покупателя (опционально)
	Nonce      int     `json:"nonce"`     // Уникальное значение для предотвращения повторных запросов
	Signature  string  `json:"signature"` // Подпись для проверки целостности данных
}

type CreateOrderResponse struct {
	Type      string `json:"type"`
	OrderID   int    `json:"orderId"`
	OrderHash string `json:"orderHash"`
	Location  string `json:"location"`
}

func createSignature(shopID int, amount float64, currency string, email string, i int, ip string, nonce int, secretKey string) string {
	message := fmt.Sprintf("%d|%f|%s|%s|%d|%s|%d", shopID, amount, currency, email, i, ip, nonce)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func createOrder(amount float64, currency string, email string, shopID int, i int, ip string, nonce int, secretKey string) (*CreateOrderResponse, error) {
	// Создание подписи
	signature := createSignature(shopID, amount, currency, email, i, ip, nonce, secretKey)

	// Подготовка данных запроса
	orderRequest := CreateOrderRequest{
		Amount:     amount,
		CurrencyID: currency,
		Email:      email,
		ShopID:     shopID,
		IP:         ip,
		I:          i,
		Nonce:      nonce,
		Signature:  signature,
	}

	requestBody, err := json.Marshal(orderRequest)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(requestBody))
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", "https://api.freekassa.ru/v1/orders/create", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var createOrderResp CreateOrderResponse
	if err := json.Unmarshal(body, &createOrderResp); err != nil {
		return nil, err
	}

	fmt.Println("resp", createOrderResp)
	return &createOrderResp, nil
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
	merchantIDToInt, _ := strconv.Atoi(merchantID)
	replenishmentIDInt, _ := strconv.Atoi(replenishmentID)
	var APIKey = os.Getenv("API_KEY")

	location, err := createOrder(amount, "RUB", email, merchantIDToInt, 44, "193.123.140.224", replenishmentIDInt, APIKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(location.Location)
	return location.Location, nil
}
