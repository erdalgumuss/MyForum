package main

import (
	"MyForum/config"
	"MyForum/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDB()

	// Create a new Gin router
	r := gin.Default()

	// Serve static files
	//r.Static("/static", "./static")

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Define routes
	routes.AuthRoutes(r)
	routes.ForumRoutes(r)
	routes.RegisterProfileRoutes(r)

	// Start the server
	r.Run(":8080")
}
