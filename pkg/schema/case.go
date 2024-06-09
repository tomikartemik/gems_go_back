package schema

import "gems_go_back/pkg/model"

type CaseInput struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type CaseInfo struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ShowCase struct {
	Id    int                `json:"id"`
	Name  string             `json:"name"`
	Price int                `json:"price"`
	Items []model.ItemWithID `json:"items"`
}
