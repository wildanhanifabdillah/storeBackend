package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/wildanhanifabdillah/storeBackend/internal/config"
	"github.com/wildanhanifabdillah/storeBackend/internal/database"
	"github.com/wildanhanifabdillah/storeBackend/internal/routes"
)

func main() {
	// 1️⃣ Load & validate config (.env)
	cfg := config.Load()
	_ = cfg // dipakai implicit (validasi + env ready)

	// 2️⃣ Init database
	db := database.InitDB()

	// 3️⃣ Init Gin
	r := gin.Default()

	// 4️⃣ Register routes
	routes.RegisterRoutes(r, db)

	// 5️⃣ Run server
	log.Println("🚀 Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
