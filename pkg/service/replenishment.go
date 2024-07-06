package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gems_go_back/pkg/repository"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
var location string
var reportTo ReportTo

type ReplenishmentService struct {
	repo repository.Replenishment
}

func NewReplenishmentService(repo repository.Replenishment) *ReplenishmentService {
	return &ReplenishmentService{repo: repo}
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

	form := url.Values{}
	form.Add("m", merchantID)
	form.Add("oa", amount)
	form.Add("o", orderID)
	form.Add("s", signature)
	form.Add("currency", currency)
	form.Add("i", description)
	form.Add("email", email)

	fmt.Println(form)
	req, err := http.NewRequest("GET", "https://pay.freekassa.com/", strings.NewReader(form.Encode()))
	fmt.Println(form.Encode())
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}
	reportToStr := resp.Header.Get("Report-To")

	err = json.Unmarshal([]byte(reportToStr), &reportTo)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	for _, endpoint := range reportTo.Endpoints {
		location = endpoint.URL
		fmt.Println("URL:", endpoint.URL)
	}

	return location, nil
}
