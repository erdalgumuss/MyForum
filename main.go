package main

import (
	"log"
	"os"

	"MyForum/config"
	"MyForum/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	config.InitOAuthConfig(clientID, clientSecret, redirectURL)

	// Initialize database
	config.InitDB()

	// Create a new Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Define routes
	routes.AuthRoutes(r)
	routes.ForumRoutes(r)
	routes.ProfileRoutes(r)

	// Start the server
	r.Run(":8080")
}
