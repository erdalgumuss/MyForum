package routes

import (
	"MyForum/controllers"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func ModeratorRoutes(r *gin.Engine) {
	// User Routes
	user := r.Group("/user")
	{
		user.POST("/request_moderator", controllers.RequestModerator)
	}

	// Moderator Routes
	moderator := r.Group("/moderator")
	moderator.Use(utils.ModeratorOnly()) // Moderatör oturum kontrolü
	{
		moderator.GET("/dashboard", controllers.ModeratorDashboard)
		moderator.GET("/pending_posts", controllers.PendingPosts)
		moderator.POST("/approve_post", controllers.ApprovePost)
		moderator.POST("/reject_post", controllers.RejectPost)
		moderator.POST("/delete_post", controllers.DeletePost) // Gönderi silme rotası
	}
}
