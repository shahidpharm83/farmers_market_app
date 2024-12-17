// models/review.go
package models

import (
	"time"
)

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id"` // The ID of the product being reviewed
	UserID    uint      `json:"user_id"`    // The ID of the user leaving the review
	Rating    int       `json:"rating"`     // Rating out of 5
	Comment   string    `json:"comment"`    // Text review/comment
	CreatedAt time.Time `json:"created_at"` // Timestamp for when the review was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp for when the review was last updated
}
