package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	MaxParticipants int       `json:"max_participants"`
	Date            time.Time `json:"date"`
	Location        string    `json:"location"`
	Status          string    `json:"status" gorm:"default:'open'"`

	VolunteerID uint      `json:"volunteer_id"`
	Volunteer   Volunteer `json:"volunteer,omitempty"`

	Registrations []ServiceRegistration `json:"registrations,omitempty"`
}

type ServiceRegistration struct {
	gorm.Model
	Status string `json:"status" gorm:"default:'registered'"`

	ServiceID uint    `json:"service_id"`
	Service   Service `json:"service,omitempty"`

	UserID uint `json:"user_id"`
	User   User `json:"user,omitempty"`
}