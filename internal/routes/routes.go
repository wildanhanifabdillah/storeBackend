package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"

	"github.com/wildanhanifabdillah/storeBackend/internal/handlers"
	"github.com/wildanhanifabdillah/storeBackend/internal/middlewares"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/health", handlers.HealthCheck())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
	admin.GET("/invoices/:order_id", handlers.AdminDownloadInvoice())

}
