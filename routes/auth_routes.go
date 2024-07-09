package routes

import (
	"MyForum/controllers"
	"MyForum/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	// Define routes
	r.GET("/", handlers.ShowIndexPage)
	r.GET("/models/user", handlers.GetUserProfile) // Yeni endpoint
	r.POST("/login", controllers.Login)
	r.POST("/register", handlers.ProcessRegister)
	r.POST("/logout", handlers.Logout)
}
