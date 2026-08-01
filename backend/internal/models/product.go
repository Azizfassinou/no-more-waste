package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	OriginalPrice float64   `json:"original_price"`
	DiscountPrice float64   `json:"discount_price"`
	Quantity      int       `json:"quantity"`
	ExpiryDate    time.Time `json:"expiry_date"`
	IsAvailable   bool      `json:"is_available"`

	MerchantID uint     `json:"merchant_id"`
	Merchant   Merchant `json:"merchant,omitempty"`
}
