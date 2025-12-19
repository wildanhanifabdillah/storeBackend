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
		log.Fatal("failed to connect database: ", err)
	}

	// 🔥 Auto Migration
	if err := db.AutoMigrate(
		&models.Admin{},
		&models.Game{},
		&models.TopupPackage{},
		&models.Transaction{},
		&models.Payment{},
	); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	// 🌱 Seed initial data
	Seed(db)

	log.Println("Database connected, migrated, and seeded successfully")

	return db
}
