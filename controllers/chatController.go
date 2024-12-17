// controllers/chatController.go
package controllers

import (
	"farmers_market_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB // Ensure this is initialized in main.go

// SendMessage handles sending a message (text, image, voice, file)
func SendMessage(c *gin.Context) {
	var message models.Message

	// Bind JSON data to the message struct
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Save the message to the database
	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

// GetMessages retrieves messages between two users
func GetMessages(c *gin.Context) {
	senderID := c.Query("sender_id")
	receiverID := c.Query("receiver_id")

	var messages []models.Message
	// Retrieve messages from the database
	if err := db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// Existing methods...
