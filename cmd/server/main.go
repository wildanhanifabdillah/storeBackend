package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/wildanhanifabdillah/storeBackend/internal/config"
	"github.com/wildanhanifabdillah/storeBackend/internal/database"
	"github.com/wildanhanifabdillah/storeBackend/internal/routes"
	"github.com/wildanhanifabdillah/storeBackend/internal/services"
)

func main() {
	// 1️⃣ Load & validate config (.env)
	cfg := config.Load()
	_ = cfg // dipakai implicit (validasi + env ready)

	// 2️⃣ Init database
	db := database.InitDB()

	// 3️⃣ Init Redis (queue email)
	// Aman walau Redis belum hidup (nanti via Docker)
	services.InitRedis()

	// 4️⃣ Start email worker (async)
	services.StartEmailWorker()

	// 5️⃣ Init Gin
	r := gin.Default()

	// 6️⃣ Register routes
	routes.RegisterRoutes(r, db)

	// 5️⃣ Run server
	log.Println("🚀 Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
