package model

type Card struct {
	Type   string `json:"type"`
	Number string `json:"number"`
	Name   string `json:"name"`
	Bank   string `json:"bank"`
}
