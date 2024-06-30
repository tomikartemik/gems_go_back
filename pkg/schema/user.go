package schema

import "gems_go_back/pkg/model"

type ShowUser struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	IsActive bool    `json:"is_active"`
	Balance  float64 `json:"balance"`
}

type InputUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserWithItems struct {
	ID       string                  `json:"id"`
	Username string                  `json:"username"`
	Email    string                  `json:"email"`
	IsActive bool                    `json:"is_active"`
	Balance  float64                 `json:"balance"`
	Items    []model.ItemOfInventory `json:"items"`
}
