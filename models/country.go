package models

import (
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name      string     `json:"name" gorm:"unique;not null"`
	Districts []District `json:"districts" gorm:"foreignKey:CountryID"` // Relationship with District
}
