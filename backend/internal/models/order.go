package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`

	ClientID uint `json:"client_id"`
	Client   User `json:"client,omitempty"`

	ProductID uint    `json:"product_id"`
	Product   Product `json:"product,omitempty"`
}
