package routes

import (
	"MyForum/controllers"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func ModeratorRoutes(r *gin.Engine) {
	// User Routes
	user := r.Group("/user")
	user.Use(utils.AuthMiddleware())
	{
		user.POST("/request_moderator", controllers.RequestModerator)
	}

	// Moderator Routes
	moderator := r.Group("/moderator")
	moderator.Use(utils.AuthMiddleware(), utils.ModeratorOnly()) // Moderator session check
	{
		moderator.GET("/dashboard", controllers.ModeratorDashboard)
		moderator.GET("/pending_posts", controllers.PendingPosts)
		moderator.POST("/approve_post", controllers.ApprovePost)
		moderator.POST("/reject_post", controllers.RejectPost)
		moderator.POST("/delete_post", controllers.DeletePost) // Delete post route
	}
}
