package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wildanhanifabdillah/storeBackend/internal/models"
)

type CreatePackageRequest struct {
	Name     string `json:"name" binding:"required"`
	Amount   int    `json:"amount" binding:"required"`
	Price    int64  `json:"price" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

type UpdatePackageRequest struct {
	Name     *string `json:"name"`
	Amount   *int    `json:"amount"`
	Price    *int64  `json:"price"`
	IsActive *bool   `json:"is_active"`
}

// GET /api/v1/admin/games/:id/packages
func AdminListPackages(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameID := c.Param("id")

		var pkgs []models.TopupPackage
		if err := db.Where("game_id = ?", gameID).Find(&pkgs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to fetch packages",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": pkgs})
	}
}

// POST /api/v1/admin/games/:id/packages
func AdminCreatePackage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameIDStr := c.Param("id")
		gameID, err := strconv.Atoi(gameIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid game id"})
			return
		}

		var req CreatePackageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
			return
		}

		pkg := models.TopupPackage{
			GameID:   uint(gameID),
			Name:     req.Name,
			Amount:   req.Amount,
			Price:    req.Price,
			IsActive: true,
		}
		if req.IsActive != nil {
			pkg.IsActive = *req.IsActive
		}

		if err := db.Create(&pkg).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create package"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": pkg})
	}
}

// PUT /api/v1/admin/packages/:id
func AdminUpdatePackage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req UpdatePackageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
			return
		}

		updates := map[string]interface{}{}
		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Amount != nil {
			updates["amount"] = *req.Amount
		}
		if req.Price != nil {
			updates["price"] = *req.Price
		}
		if req.IsActive != nil {
			updates["is_active"] = *req.IsActive
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no fields to update"})
			return
		}

		if err := db.Model(&models.TopupPackage{}).
			Where("id = ?", id).
			Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update package"})
			return
		}

		var pkg models.TopupPackage
		if err := db.First(&pkg, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load package"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": pkg})
	}
}

// DELETE /api/v1/admin/packages/:id
func AdminDeletePackage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := db.Delete(&models.TopupPackage{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete package"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "package deleted"})
	}
}

