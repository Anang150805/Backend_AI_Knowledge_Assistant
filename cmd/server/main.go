package main

import (
	"log"
	"os"

	"backend-AI-Knowledge-Assistant/internal/db"
	"backend-AI-Knowledge-Assistant/internal/routes"

	"github.com/gin-contrib/cors" // 🔥 tambahan
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 🔌 Koneksi database
	conn, err := db.NewPostgres()
	if err != nil {
		log.Fatal("failed to connect db:", err)
	}
	defer conn.Close()

	// 🚀 Init Gin
	r := gin.Default()

	// 🔥 CORS CONFIG (FIX ERROR OPTIONS 404)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173", // frontend Vite
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		AllowCredentials: true,
	}))

	// 📌 Routes
	routes.SetupRoutes(r, conn)

	// 🌐 Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server running on port", port)

	// ▶️ Run server
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}