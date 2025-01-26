package model

type OwnReplenishment struct {
	ID         int     `json:"id" db:"id" gorm:"autoIncrement"`
	UserId     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	ReceiptURL string  `json:"receipt_url"`
	Status     string  `json:"status"`
}

type OwnReplenishmentsResponse struct {
	ID         int     `json:"id" db:"id" gorm:"autoIncrement"`
	UserId     string  `json:"user_id"`
	Username   string  `json:"username"`
	Amount     float64 `json:"amount"`
	ReceiptURL string  `json:"receipt_url"`
	Status     string  `json:"status"`
}

type OwnReplenishmentOutput struct {
	Replenishments []OwnReplenishmentsResponse `json:"replenishments"`
	PagesCount     int                         `json:"pages_count"`
}
