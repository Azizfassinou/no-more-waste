package models

import (
	"time"

	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	CompanyName    string    `json:"company_name"`
	SiretNumber    string    `json:"siret_number" gorm:"unique"`
	CompanyAddress string    `json:"company_address"`
	IsApproved     bool      `json:"is_approved" gorm:"default:false"`
	IsActive       bool      `json:"is_active"`
	RenewalDate    time.Time `json:"renewal_date"`

	UserID uint `json:"user_id"`
	User   User `json:"user,omitempty"`
}
