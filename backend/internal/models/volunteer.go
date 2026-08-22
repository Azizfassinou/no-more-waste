package models

import (
	"gorm.io/gorm"
)

type Volunteer struct {
	gorm.Model
	ZoneArea     string `json:"zone_area"`
	Availability string `json:"availability"`
	Vehicle      bool   `json:"vehicle"`

	UserID uint    `json:"user_id"`
	User   User    `json:"user,omitempty"`
	Skills []Skill `json:"skills,omitempty" gorm:"many2many:volunteer_skills;"`
}