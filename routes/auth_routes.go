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
	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)
	r.GET("/auth/github/login", handlers.GitHubLogin)
    r.GET("/auth/github/callback", handlers.GitHubCallback)
}
