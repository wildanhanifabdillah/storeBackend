package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/wildanhanifabdillah/storeBackend/internal/models"
)

func Seed(db *gorm.DB) {
	seedAdmin(db)
	seedGamesAndPackages(db)
}

func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		log.Println("Admin already seeded")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	admin := models.Admin{
		Name:         "Super Admin",
		Email:        "admin@mail.com",
		PasswordHash: string(password),
	}

	db.Create(&admin)
	log.Println("Admin seeded")
}

func seedGamesAndPackages(db *gorm.DB) {
	var count int64
	db.Model(&models.Game{}).Count(&count)
	if count > 0 {
		log.Println("Games already seeded")
		return
	}

	// ===== GAME 1 =====
	ml := models.Game{
		Name: "Mobile Legends",
		Code: "ml",
	}
	db.Create(&ml)

	mlPackages := []models.TopupPackage{
		{Name: "86 Diamonds", Amount: 86, Price: 20000, GameID: ml.ID},
		{Name: "172 Diamonds", Amount: 172, Price: 40000, GameID: ml.ID},
	}

	db.Create(&mlPackages)

	// ===== GAME 2 =====
	ff := models.Game{
		Name: "Free Fire",
		Code: "ff",
	}
	db.Create(&ff)

	ffPackages := []models.TopupPackage{
		{Name: "140 Diamonds", Amount: 140, Price: 20000, GameID: ff.ID},
		{Name: "355 Diamonds", Amount: 355, Price: 50000, GameID: ff.ID},
	}

	db.Create(&ffPackages)

	log.Println("Games & packages seeded")
}
