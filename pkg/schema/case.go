package schema

import "gems_go_back/pkg/model"

type CaseInput struct {
	Name      string `json:"name"`
	Price     int    `json:"price"`
	PhotoLink string `json:"photo_link"`
	Color     string `json:"color"`
}

type CaseInfo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	PhotoLink string `json:"photo_link"`
	Color     string `json:"color"`
}

type ShowCase struct {
	Id        int                `json:"id"`
	Name      string             `json:"name"`
	Price     int                `json:"price"`
	PhotoLink string             `json:"photo_link"`
	Items     []model.ItemWithID `json:"items"`
}
