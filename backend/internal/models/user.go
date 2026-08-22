package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	Address        string     `json:"address"`
	ResetCode      string     `json:"-"`
	ResetExpiresAt *time.Time `json:"-"`
}