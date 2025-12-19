package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wildanhanifabdillah/storeBackend/internal/handlers"
	"github.com/wildanhanifabdillah/storeBackend/internal/middlewares"
	"github.com/wildanhanifabdillah/storeBackend/internal/services"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	// 🔥 TEST EMAIL (LOCAL ONLY)
	r.GET("/test-email", func(c *gin.Context) {
		err := services.SendPaymentSuccessEmail(
			"wildanhanifabdillah27@gmail.com",
			"TEST-LOCAL",
			20000,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "email sent"})
	})

	api := r.Group("/api/v1")

	// Public
	api.GET("/games", handlers.GetGames(db))
	api.GET("/games/:id/packages", handlers.GetPackages(db))
	api.POST("/checkout", handlers.Checkout(db))
	api.GET("/transactions/:order_id", handlers.GetTransactionStatus(db))

	// Payment callback
	api.POST("/payments/midtrans/callback", handlers.MidtransCallback(db))

	// Admin
	admin := api.Group("/admin")
	admin.POST("/login", handlers.AdminLogin(db))

	admin.Use(middlewares.AuthMiddleware())
	admin.GET("/transactions", handlers.AdminGetTransactions(db))
	admin.POST("/games", handlers.CreateGame(db))
}
