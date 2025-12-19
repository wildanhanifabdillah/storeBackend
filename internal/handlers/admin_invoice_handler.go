package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/admin/invoices/:order_id
func AdminDownloadInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		path := "invoices/" + orderID + ".pdf"

		if _, err := os.Stat(path); os.IsNotExist(err) {
			c.JSON(404, gin.H{"message": "invoice not found"})
			return
		}

		c.File(path)
	}
}
