package models

import (
	"time"

	"gorm.io/gorm"
)

type District struct {
	gorm.Model
	Name      string    `json:"name" gorm:"unique;not null"`
	CountryID uint      `json:"country_id"`
	Country   Country   `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
