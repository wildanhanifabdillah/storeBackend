package handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wildanhanifabdillah/storeBackend/internal/models"
	"github.com/wildanhanifabdillah/storeBackend/internal/services"
)

type MidtransCallbackRequest struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
}

// POST /api/v1/payments/midtrans/callback
func MidtransCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MidtransCallbackRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid payload",
			})
			return
		}

		// ===============================
		// 1️⃣ VERIFY SIGNATURE MIDTRANS
		// ===============================
		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		rawSignature := req.OrderID + req.GrossAmount + req.TransactionStatus + serverKey

		hash := sha512.Sum512([]byte(rawSignature))
		expectedSignature := hex.EncodeToString(hash[:])

		if req.SignatureKey != expectedSignature {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid signature",
			})
			return
		}

		// ===============================
		// 2️⃣ FIND TRANSACTION
		// ===============================
		var tx models.Transaction
		if err := db.Where("order_id = ?", req.OrderID).
			First(&tx).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "transaction not found",
			})
			return
		}

		// ===============================
		// 3️⃣ MAP STATUS MIDTRANS
		// ===============================
		var newStatus string
		switch req.TransactionStatus {
		case "settlement", "capture":
			newStatus = "PAID"
		case "expire", "cancel", "deny":
			newStatus = "FAILED"
		default:
			newStatus = "PENDING"
		}

		// ===============================
		// 4️⃣ UPDATE TRANSACTION
		// ===============================
		tx.Status = newStatus
		if err := db.Save(&tx).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to update transaction",
			})
			return
		}

		// ===============================
		// 5️⃣ UPDATE PAYMENT
		// ===============================
		if err := db.Model(&models.Payment{}).
			Where("transaction_id = ?", tx.ID).
			Updates(map[string]interface{}{
				"payment_status": req.TransactionStatus,
				"payment_type":   req.PaymentType,
			}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to update payment",
			})
			return
		}

		// ===============================
		// 6️⃣ SEND EMAIL IF PAID
		// ===============================
		if newStatus == "PAID" {
			// Email failure SHOULD NOT break payment flow
			_ = services.SendPaymentSuccessEmail(
				tx.Email,
				tx.OrderID,
				tx.TotalAmount,
			)
		}

		// ===============================
		// 7️⃣ RESPONSE TO MIDTRANS
		// ===============================
		c.JSON(http.StatusOK, gin.H{
			"message": "callback processed successfully",
		})
	}
}
