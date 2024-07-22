package main

import (
	"log"
	"os"
	"path/filepath"

	"MyForum/config"
	"MyForum/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	// Session middleware
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Serve static files
	r.Static("/static", "../frontend/static")
	//r.Static("/uploads", "../frontend/uploads")
	r.Static("/uploads", filepath.Join("frontend", "uploads"))

	// Load HTML templates
	r.LoadHTMLGlob("../frontend/templates/*")

	// Define routes
	routes.AuthRoutes(r)
	routes.ForumRoutes(r)
	routes.ProfileRoutes(r)
	routes.AdminRoutes(r)
	routes.ModeratorRoutes(r)

	// Start the server
	r.Run(":8080")
}
