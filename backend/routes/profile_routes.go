package routes

import (
	"MyForum/controllers"
	"MyForum/handlers"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(utils.AuthMiddleware())
	{
		protected.GET("/models/user", controllers.GetUserProfile) // Yeni endpoint

		// protected.GET("/profile", handlers.ShowProfilePage)
		protected.GET("/profile", handlers.ProfileView)
		protected.GET("/user/:id/posts", handlers.GetUserPosts)
		protected.GET("/user/:id/likes", handlers.GetUserLikes)
		protected.GET("/user/:id/comments", handlers.GetUserComments)

		protected.GET("/inbox", handlers.ShowInbox)
		protected.GET("/get-messages", handlers.GetMessages)
		protected.POST("/send-message", handlers.SendMessages)

		// router.PUT("/profile", utils.AuthMiddleware(), handlers.ProfileUpdate)
		// router.POST("/profile/change-password", utils.AuthMiddleware(), handlers.ChangePassword)

		protected.GET("/profile/", controllers.GetUserProfile)
		// router.POST("/profile/change-password", utils.AuthRequired(), controllers.ChangePassword)
	}
}
