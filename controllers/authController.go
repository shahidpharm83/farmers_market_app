package controllers

import (
	"farmers_market_backend/config"
	"farmers_market_backend/middleware"
	"farmers_market_backend/models"
	"farmers_market_backend/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Register registers a new user
//
// This function takes a JSON payload from the request body, validates
// the input, hashes the password, and creates a new user in the database.
//
// The input struct contains three fields: username, email, and password.
// The username must be alphanumeric and between 3 and 255 characters long.
// The email must be a valid email address and between 6 and 255 characters long.
// The password must be at least 6 characters long and at most 255 characters long.
//
// If the input is invalid, the function returns a 400 Bad Request response
// with an error message.
//
// If the password hashing fails, the function returns a 500 Internal Server
// Error response with an error message.
//
// If the user is saved to the database successfully, the function returns
// a 201 Created response with a success message.
func Register(c *gin.Context) {
	// Input struct to store the user's registration data
	var input struct {
		Username string `json:"username" binding:"required,alphanum,max=255,min=3"`
		Email    string `json:"email" binding:"required,email,max=255,min=6"`
		Password string `json:"password" binding:"required,min=6,max=255"`
	}

	// Bind input data from the request body to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EOF"})
		return
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user with the input data
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Save the user to the database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Return a success message if the user is registered successfully
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	config.DB.Where("email = ?", input.Email).First(&user)

	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// UpdateUser updates an existing user by ID
func UpdateUser(c *gin.Context) {
	userID := c.Param("userID")
	var updatedUser models.User

	// Bind JSON data to the updatedUser struct
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the existing user
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	// Update user fields
	user.Name = updatedUser.Name
	user.ImageURL = updatedUser.ImageURL // Update the image URL if provided
	user.IsSeller = updatedUser.IsSeller // Update the isSeller status

	// Save the updated user
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUsers retrieves all users
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Preload("District").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser fetches a user by ID. Requires either admin rights or self-ownership of the profile.
func GetUser(c *gin.Context) {
	// Parse user ID from URL or middleware
	userIDParam := c.Param("id")
	requestedUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user ID from the context set by middleware (e.g., after token validation)
	authUserID, exists := c.Get("user_id") // Ensure middleware sets user_id in context
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// Check if the requested user is either the authenticated user or the user is an admin
	if uint(authUserID.(uint)) != uint(requestedUserID) && !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	// Retrieve user details from database
	var user models.User
	if err := config.DB.Preload("District").First(&user, requestedUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with user data
	c.JSON(http.StatusOK, gin.H{
		"id":               user.ID,
		"name":             user.Name,
		"email":            user.Email,
		"username":         user.Username,
		"image_url":        user.ImageURL,
		"is_seller":        user.IsSeller,
		"district_id":      user.DistrictID,
		"district":         user.District,
		"delivery_address": user.DeliveryAddress,
		"mobile_number":    user.MobileNumber,
		"is_admin":         user.IsAdmin,
		"created_at":       user.CreatedAt,
		"updated_at":       user.UpdatedAt,
	})
}

// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	userID := c.Param("userID")

	// Delete the user by ID
	if err := db.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Example function to check if user is admin
func IsAdmin(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_admin": user.IsAdmin})
}

// RefreshToken handles the refresh token request
func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Refresh token: %v", input.RefreshToken)

	// Validate the refresh token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Use the JWTSecretKey variable instead of the hardcoded secret key
		return middleware.JWTSecretKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	// Extract user ID from claims
	userID := uint(claims["sub"].(float64)) // Ensure "sub" is set correctly in your token

	// Generate a new access token
	newAccessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
