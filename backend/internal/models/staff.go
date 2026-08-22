package models

import (
	"gorm.io/gorm"
)

type Staff struct {
	gorm.Model
	Department string  `json:"department"`
	JobTitle   string  `json:"job_title"`
	Salary     float64 `json:"salary"`

	UserID uint `json:"user_id"`
	User   User `json:"user,omitempty"`
}