package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/wildanhanifabdillah/storeBackend/internal/database"
	"github.com/wildanhanifabdillah/storeBackend/internal/routes"
)

func main() {
	_ = godotenv.Load()

	db := database.InitDB()

	r := gin.Default()
	routes.RegisterRoutes(r, db)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
