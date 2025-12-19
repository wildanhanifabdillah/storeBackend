package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/wildanhanifabdillah/storeBackend/internal/models"
)

func InitDB() *gorm.DB {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// 🔥 Auto Migration
	err = db.AutoMigrate(
		&models.Admin{},
		&models.Game{},
		&models.TopupPackage{},
		&models.Transaction{},
		&models.Payment{},
	)
	if err != nil {
		log.Fatal("failed to migrate database")
	}

	return db
}
