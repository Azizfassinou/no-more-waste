package models

import (
	"time"

	"gorm.io/gorm"
)

type CollectionMission struct {
	gorm.Model
	Status     string    `json:"status" gorm:"default:'assigned'"`
	PickupDate time.Time `json:"pickup_date"`

	MerchantID uint     `json:"merchant_id"`
	Merchant   Merchant `json:"merchant,omitempty"`

	VolunteerID uint      `json:"volunteer_id"`
	Volunteer   Volunteer `json:"volunteer,omitempty"`

	ProductID uint    `json:"product_id"`
	Product   Product `json:"product,omitempty"`
}