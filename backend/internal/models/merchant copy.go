package models

import (
	"time"

	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email" gorm:"unique"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	RenewalDate time.Time `json:"renewal_date"`
	IsActive    bool      `json:"is_active"`
}
