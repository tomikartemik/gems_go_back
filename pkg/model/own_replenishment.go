package model

type OwnReplenishment struct {
	ID         int     `json:"id" db:"id" gorm:"autoIncrement"`
	UserId     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	ReceiptURL string  `json:"receipt_url"`
	Status     string  `json:"status"`
}

type OwnReplenishmentOutput struct {
	Replenishments []OwnReplenishment `json:"replenishments"`
	PagesCount     int                `json:"pages_count"`
}
