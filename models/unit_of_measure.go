// models/unit_of_measure.go
package models

import (
	"time"
)

// UnitOfMeasure struct
type UnitOfMeasure struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"unique;not null"`
	Abbreviation string    `json:"abbreviation"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
