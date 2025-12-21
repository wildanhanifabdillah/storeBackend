package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/wildanhanifabdillah/storeBackend/internal/config"
	"github.com/wildanhanifabdillah/storeBackend/internal/database"
	"github.com/wildanhanifabdillah/storeBackend/internal/routes"
	"github.com/wildanhanifabdillah/storeBackend/internal/services"
)

func main() {
	// 1?,???? Load & validate config (.env)
	cfg := config.Load()
	_ = cfg // dipakai implicit (validasi + env ready)

	// 2?,???? Init database
	db := database.InitDB()

	// 3?,???? Init Redis (queue email)
	// Aman walau Redis belum hidup (nanti via Docker)
	services.InitRedis()

	// 4?,???? Start email worker (async)
	services.StartEmailWorker()

	// 5?,???? Init Gin
	r := gin.Default()

	// Allow frontend origins
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://whastore.my.id",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 6?,???? Register routes
	routes.RegisterRoutes(r, db)

	// 5?,???? Run server
	log.Println("dYs? Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
