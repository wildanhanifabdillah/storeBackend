package handlers

import (
	"net/http"

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
		c.JSON(http.StatusOK, gin.H{"message": "AdminGetTransactions OK"})
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
