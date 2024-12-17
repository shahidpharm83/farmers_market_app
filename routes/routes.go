package routes

import (
	"farmers_market_backend/controllers"
	"farmers_market_backend/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Authentication routes
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
	api.POST("/refresh-token", controllers.RefreshToken)

	userGroup := router.Group("/api/users")
	{
		userGroup.GET("/", controllers.GetUsers)                                  // Route for getting all users
		userGroup.GET("/:id", middleware.IsItMyProfile(), controllers.GetUser)    // Route for getting a specific user by ID
		userGroup.PUT("/:id", middleware.IsItMyProfile(), controllers.UpdateUser) // Route for updating user profile
		userGroup.DELETE("/:id", middleware.IsAdmin(), controllers.DeleteUser)    // Route for updating user profile

	}

	userGroup.Use(middleware.AuthMiddleware())

	// Country routes
	countryRoutes := router.Group("/api/countries")
	{
		countryRoutes.GET("/countries", controllers.GetCountries)
		countryRoutes.GET("/countries/:id", controllers.GetCountry)
		countryRoutes.POST("/countries", controllers.CreateCountry)
		countryRoutes.PUT("/countries/:id", controllers.UpdateCountry)
		countryRoutes.DELETE("/countries/:id", controllers.DeleteCountry)
	}
	countryRoutes.Use(middleware.AuthMiddleware())

	// District routes
	districtRoutes := router.Group("/api/districts")
	{
		districtRoutes.GET("/districts", controllers.GetDistricts)
		districtRoutes.GET("/districts/:id", controllers.GetDistrict)
		districtRoutes.POST("/districts", controllers.CreateDistrict)
		districtRoutes.PUT("/districts/:id", middleware.IsAdmin(), controllers.UpdateDistrict)
		districtRoutes.DELETE("/districts/:id", middleware.AuthMiddleware(), middleware.IsAdmin(), controllers.DeleteDistrict)
	}
	districtRoutes.Use(middleware.AuthMiddleware())

	// Product routes
	productRoutes := router.Group("/api/products")
	{
		productRoutes.GET("/:productID", controllers.GetProduct)                                            // Get a single product
		productRoutes.GET("/", controllers.GetAllProducts)                                                  // Get all products
		productRoutes.POST("/", controllers.CreateProduct)                                                  // Create a new product
		productRoutes.PUT("/:productID", middleware.IsSellerOfTheItem(), controllers.UpdateProduct)         // Update an existing product
		productRoutes.DELETE("/:productID", middleware.IsCustomerForThisOrder(), controllers.DeleteProduct) // Delete a product
	}
	productRoutes.Use(middleware.AuthMiddleware())

	// Order routes

	orderRoutes := router.Group("/api/orders")
	{
		orderRoutes.POST("/", controllers.PlaceOrder)                                                  // Create a new order
		orderRoutes.GET("/:orderID", middleware.IsCustomerForThisOrder(), controllers.GetOrderDetails) // Get order details
		orderRoutes.PUT("/:orderID", middleware.IsCustomerForThisOrder(), controllers.UpdateOrder)     // Update an existing order
		orderRoutes.DELETE("/:orderID", middleware.IsCustomerForThisOrder(), controllers.DeleteOrder)  // Delete an order
	}
	orderRoutes.Use(middleware.AuthMiddleware())

	// Wallet routes
	walletRoutes := router.Group("/api/wallet")
	{
		walletRoutes.POST("/wallet/topup", controllers.TopUpWallet)
		walletRoutes.GET("/wallet/balance", controllers.GetWalletBalance)

	}
	walletRoutes.Use(middleware.AuthMiddleware())

	// Chat routes
	chatRoutes := router.Group("/api/chat")
	{
		chatRoutes.POST("/send", controllers.SendMessage)    // Send a message
		chatRoutes.GET("/messages", controllers.GetMessages) // Get messages between users
	}
	chatRoutes.Use(middleware.AuthMiddleware())

	// Reviews routes
	reviewRoutes := router.Group("/api/review")
	{
		reviewRoutes.GET("/reviews/:productID", controllers.GetReviews)  // Get all reviews for a product
		reviewRoutes.POST("/reviews", controllers.CreateReview)          // Create a review
		reviewRoutes.PUT("/reviews/:reviewID", controllers.UpdateReview) // Update a review
		reviewRoutes.DELETE("/reviews/:reviewID", controllers.DeleteReview)
	}
	reviewRoutes.Use(middleware.AuthMiddleware())

	// Category Routes
	categoryRoutes := router.Group("/categories")
	{
		categoryRoutes.POST("/", middleware.IsAdmin(), controllers.CreateCategory)              // Create a new category
		categoryRoutes.GET("/", controllers.GetCategories)                                      // Get all categories
		categoryRoutes.PUT("/:categoryID", middleware.IsAdmin(), controllers.UpdateCategory)    // Update a category
		categoryRoutes.DELETE("/:categoryID", middleware.IsAdmin(), controllers.DeleteCategory) // Delete a category
	}
	categoryRoutes.Use(middleware.AuthMiddleware())

	// Unit of Measure Routes
	unitRoutes := router.Group("/units")
	{
		unitRoutes.POST("/", middleware.IsAdmin(), controllers.CreateUnitOfMeasure)          // Create a new unit of measure
		unitRoutes.GET("/", controllers.GetUnitsOfMeasure)                                   // Get all units of measure
		unitRoutes.PUT("/:unitID", middleware.IsAdmin(), controllers.UpdateUnitOfMeasure)    // Update a unit of measure
		unitRoutes.DELETE("/:unitID", middleware.IsAdmin(), controllers.DeleteUnitOfMeasure) // Delete a unit of measure
	}
	unitRoutes.Use(middleware.AuthMiddleware())

}
