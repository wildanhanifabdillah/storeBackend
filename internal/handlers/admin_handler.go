package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/wildanhanifabdillah/storeBackend/internal/models"
	"github.com/wildanhanifabdillah/storeBackend/internal/services"
)

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AdminLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		var admin models.Admin
		if err := db.Where("email = ?", req.Email).First(&admin).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(admin.PasswordHash),
			[]byte(req.Password),
		); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
			return
		}

		token, err := services.GenerateAdminToken(admin.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

func AdminGetTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status") // optional: PENDING/PAID/FAILED

		var txs []models.Transaction
		query := db.Preload("Game").Preload("Package").Preload("Payment").Order("created_at DESC")
		if status != "" {
			query = query.Where("status = ?", status)
		}

		if err := query.Find(&txs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to fetch transactions",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": txs,
		})
	}
}

// GET /api/v1/admin/games
func AdminListGames(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var games []models.Game
		if err := db.Find(&games).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to fetch games",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": games})
	}
}

func CreateGame(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		name := c.PostForm("name")
		code := c.PostForm("code")

		if name == "" || code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "name and code are required",
			})
			return
		}

		// cek unik code
		var count int64
		db.Model(&models.Game{}).
			Where("code = ?", code).
			Count(&count)

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "game code already exists",
			})
			return
		}

		fileHeader, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "image is required",
			})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed to open image",
			})
			return
		}
		defer file.Close()

		imageURL, err := services.UploadImage(file, "games")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to upload image",
			})
			return
		}

		game := models.Game{
			Name:     name,
			Code:     code,
			ImageURL: imageURL,
			IsActive: true,
		}

		if err := db.Create(&game).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create game",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": game,
		})
	}
}

// PUT /api/v1/admin/games/:id
func UpdateGame(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var game models.Game
		if err := db.First(&game, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "game not found"})
			return
		}

		name := c.PostForm("name")
		code := c.PostForm("code")
		isActiveStr := c.PostForm("is_active")

		// cek unik code jika diubah
		if code != "" && code != game.Code {
			var count int64
			db.Model(&models.Game{}).
				Where("code = ? AND id <> ?", code, id).
				Count(&count)
			if count > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "game code already exists",
				})
				return
			}
			game.Code = code
		}

		if name != "" {
			game.Name = name
		}

		if isActiveStr != "" {
			if parsed, err := strconv.ParseBool(isActiveStr); err == nil {
				game.IsActive = parsed
			}
		}

		// Optional: update image jika ada file
		if fileHeader, err := c.FormFile("image"); err == nil && fileHeader != nil {
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "failed to open image",
				})
				return
			}
			defer file.Close()

			imageURL, err := services.UploadImage(file, "games")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to upload image",
				})
				return
			}
			game.ImageURL = imageURL
		}

		if err := db.Save(&game).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update game"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": game})
	}
}

// DELETE /api/v1/admin/games/:id
func DeleteGame(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// hapus packages terlebih dahulu (optional: soft delete instead)
		if err := db.Where("game_id = ?", id).Delete(&models.TopupPackage{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete packages"})
			return
		}

		if err := db.Delete(&models.Game{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete game"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "game deleted"})
	}
}
