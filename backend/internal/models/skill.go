package models

import "gorm.io/gorm"

type Skill struct {
	gorm.Model
	Name        string      `json:"name" gorm:"unique"`
	Category    string      `json:"category"`
	Volunteers  []Volunteer `json:"volunteers,omitempty" gorm:"many2many:volunteer_skills;"`
}