package model

type CaseRecord struct {
	ID     int `json:"id" db:"id" gorm:"autoIncrement"`
	CaseID int `json:"case_id"`
}
