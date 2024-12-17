// models/message.go
package models

import (
	"time"
)

// Message struct
type Message struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    SenderID     uint      `json:"sender_id"`
    ReceiverID   uint      `json:"receiver_id"`
    Content      string    `json:"content"`            // Text content of the message
    ImageURL     string    `json:"image_url"`         // URL for an image attachment
    VoiceURL     string    `json:"voice_url"`         // URL for a voice message
    FileURL      string    `json:"file_url"`          // URL for file attachment
    CreatedAt    time.Time `json:"created_at"`        // Timestamp for message creation
}
