package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wildanhanifabdillah/storeBackend/internal/models"
	"github.com/wildanhanifabdillah/storeBackend/internal/services"
)

// GET /api/v1/games
func GetGames(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var games []models.Game

		if err := db.Where("is_active = ?", true).
			Find(&games).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to fetch games",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": games,
		})
	}
}

// GET /api/v1/games/:id/packages
func GetPackages(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameID := c.Param("id")

		var packages []models.TopupPackage

		if err := db.
			Where("game_id = ? AND is_active = ?", gameID, true).
			Find(&packages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to fetch packages",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": packages,
		})
	}
}

type CheckoutRequest struct {
	GameID     uint   `json:"game_id"`
	PackageID  uint   `json:"package_id"`
	GameUserID string `json:"game_user_id"`
	Email      string `json:"email"`
}

// POST /api/v1/checkout
func Checkout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CheckoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}

		// 1️⃣ Pastikan game aktif
		var game models.Game
		if err := db.Where("id = ? AND is_active = ?", req.GameID, true).
			First(&game).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "game not found",
			})
			return
		}

		// 2️⃣ Pastikan package valid
		var pkg models.TopupPackage
		if err := db.Where(
			"id = ? AND game_id = ? AND is_active = ?",
			req.PackageID, req.GameID, true,
		).First(&pkg).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "package not found",
			})
			return
		}

		// 3️⃣ Generate order_id
		orderID := services.GenerateOrderID()

		// 4️⃣ Simpan transaksi
		tx := models.Transaction{
			OrderID:     orderID,
			GameID:      game.ID,
			PackageID:   pkg.ID,
			GameUserID:  req.GameUserID,
			Email:       req.Email,
			TotalAmount: pkg.Price,
			Status:      "PENDING",
		}

		if err := db.Create(&tx).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create transaction",
			})
			return
		}

		// 5️⃣ CREATE SNAP TOKEN (MIDTRANS)
		snap, err := services.CreateSnap(orderID, pkg.Price, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create payment",
			})
			return
		}

		// 6️⃣ SIMPAN PAYMENT
		payment := models.Payment{
			TransactionID:  tx.ID,
			PaymentGateway: "midtrans",
			PaymentToken:   snap.Token,
			PaymentStatus:  "pending",
		}

		if err := db.Create(&payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to save payment",
			})
			return
		}

		// 7️⃣ RESPONSE KE FRONTEND
		c.JSON(http.StatusOK, gin.H{
			"order_id":     orderID,
			"total_amount": pkg.Price,
			"snap_token":   snap.Token,
			"redirect_url": snap.RedirectURL,
		})
	}
}

// GET /api/v1/transactions/:order_id
func GetTransactionStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")

		var tx models.Transaction
		if err := db.Where("order_id = ?", orderID).
			First(&tx).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "transaction not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"order_id":     tx.OrderID,
			"status":       tx.Status,
			"total_amount": tx.TotalAmount,
		})
	}
}
