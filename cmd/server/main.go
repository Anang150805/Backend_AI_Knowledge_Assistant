package main

import (
	"log"
	"os"

	"backend-AI-Knowledge-Assistant/internal/db"
	"backend-AI-Knowledge-Assistant/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	conn, err := db.NewPostgres()
	if err != nil {
		log.Fatal("failed to connect db:", err)
	}
	defer conn.Close()

	r := gin.Default()
	routes.SetupRoutes(r, conn)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}