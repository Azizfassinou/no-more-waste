package models

import (
	"time"

	"gorm.io/gorm"
)

type DistributionRound struct {
	gorm.Model
	Date   time.Time `json:"date"`
	Status string    `json:"status" gorm:"default:'planned'"`
	Notes  string    `json:"notes"`

	VolunteerID uint      `json:"volunteer_id"`
	Volunteer   Volunteer `json:"volunteer,omitempty"`

	Deliveries []Delivery `json:"deliveries,omitempty"`
}

type Delivery struct {
	gorm.Model
	RecipientName    string `json:"recipient_name"`
	RecipientAddress string `json:"recipient_address"`
	RecipientType    string `json:"recipient_type"`
	Quantity         int    `json:"quantity"`
	Status           string `json:"status" gorm:"default:'pending'"`

	DistributionRoundID uint              `json:"distribution_round_id"`
	DistributionRound   DistributionRound `json:"-"`

	ProductID uint    `json:"product_id"`
	Product   Product `json:"product,omitempty"`
}