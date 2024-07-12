package main

import (
	"MyForum/config"
	"MyForum/routes"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDB()
	utils.InitGoogleOAuth("YOUR_GOOGLE_CLIENT_ID", "YOUR_GOOGLE_CLIENT_SECRET", "http://localhost:8080/auth/google/callback")

	// Create a new Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Define routes
	routes.AuthRoutes(r)
	routes.ForumRoutes(r)
	routes.ProfileRoutes(r)

	// Start the server
	r.Run(":8080")
}
