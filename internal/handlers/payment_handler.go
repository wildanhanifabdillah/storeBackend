package handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
			return
		}

		log.Printf("Midtrans callback: order_id=%s status=%s",
			req.OrderID, req.TransactionStatus)

		// ===============================
		// 1️⃣ VERIFY SIGNATURE MIDTRANS
		// ===============================
		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		rawSignature := req.OrderID + req.GrossAmount + req.TransactionStatus + serverKey

		hash := sha512.Sum512([]byte(rawSignature))
		expectedSignature := hex.EncodeToString(hash[:])

		if req.SignatureKey != expectedSignature {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid signature"})
			return
		}

		// ===============================
		// 2️⃣ FIND TRANSACTION
		// ===============================
		var tx models.Transaction
		if err := db.Where("order_id = ?", req.OrderID).
			First(&tx).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
			return
		}

		// ===============================
		// 2️⃣.5️⃣ IDEMPOTENCY
		// ===============================
		if tx.Status == "PAID" {
			c.JSON(http.StatusOK, gin.H{"message": "transaction already processed"})
			return
		}

		// ===============================
		// 3️⃣ VALIDATE AMOUNT
		// ===============================
		if req.GrossAmount != fmt.Sprintf("%d", tx.TotalAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "amount mismatch"})
			return
		}

		// ===============================
		// 4️⃣ MAP STATUS MIDTRANS
		// ===============================
		var newStatus string
		switch req.TransactionStatus {
		case "settlement":
			newStatus = "PAID"
		case "capture":
			if req.PaymentType == "credit_card" {
				newStatus = "PAID"
			} else {
				newStatus = "PENDING"
			}
		case "expire", "cancel", "deny":
			newStatus = "FAILED"
		default:
			newStatus = "PENDING"
		}

		// ===============================
		// 5️⃣ UPDATE TRANSACTION
		// ===============================
		tx.Status = newStatus
		if err := db.Save(&tx).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update transaction"})
			return
		}

		// ===============================
		// 6️⃣ UPDATE PAYMENT
		// ===============================
		if err := db.Model(&models.Payment{}).
			Where("transaction_id = ?", tx.ID).
			Updates(map[string]interface{}{
				"payment_status": req.TransactionStatus,
				"payment_type":   req.PaymentType,
			}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update payment"})
			return
		}

		// ===============================
		// 7️⃣ LOAD GAME & PACKAGE (JOIN DB)
		// ===============================
		var game models.Game
		if err := db.First(&game, tx.GameID).Error; err != nil {
			log.Println("failed load game:", err)
		}

		var pkg models.TopupPackage
		if err := db.First(&pkg, tx.PackageID).Error; err != nil {
			log.Println("failed load package:", err)
		}

		// ===============================
		// 8️⃣ GENERATE INVOICE + SEND EMAIL
		// ===============================
		if newStatus == "PAID" {

			// 🔥 Generate Invoice PDF (REAL DATA)
			invoicePath, err := services.GenerateInvoicePDF(
				services.InvoiceData{
					OrderID:     tx.OrderID,
					Email:       tx.Email,
					GameName:    game.Name,
					PackageName: pkg.Name,
					Amount:      tx.TotalAmount,
					PaidAt:      time.Now(),
				},
			)
			if err != nil {
				log.Println("failed generate invoice:", err)
			}

			// 📧 Send Email + ATTACH Invoice
			if err := services.SendPaymentSuccessEmailWithInvoice(
				tx.Email,
				tx.OrderID,
				tx.TotalAmount,
				invoicePath,
			); err != nil {
				log.Println("failed send email:", err)
			}
		}

		// ===============================
		// 9️⃣ RESPONSE TO MIDTRANS
		// ===============================
		c.JSON(http.StatusOK, gin.H{"message": "callback processed successfully"})
	}
}
