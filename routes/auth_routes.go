package routes

import (
	"MyForum/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	// Define routes
	r.GET("/", handlers.ShowIndexPage)
	r.POST("/login", handlers.ProcessLogin)
	r.POST("/register", handlers.ProcessRegister)
	r.POST("/logout", handlers.Logout)
}
