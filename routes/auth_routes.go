package routes

import (
	"MyForum/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	// Define routes
	r.GET("/", handlers.ShowIndexPage)
	// r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.ProcessLogin)
	//r.GET("/register", handlers.ShowRegisterPage)
	r.POST("/register", handlers.ProcessRegister)
	r.POST("/logout", handlers.Logout)
}
