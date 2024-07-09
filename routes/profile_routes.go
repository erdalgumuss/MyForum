package routes

import (
	"MyForum/handlers"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(utils.AuthMiddleware())
	{
		protected.GET("/models/user", handlers.GetUserProfile) // Yeni endpoint

		// protected.GET("/profile", handlers.ShowProfilePage)
		 protected.GET("/profile", handlers.ProfileView)

	}
}

// router.PUT("/profile", utils.AuthMiddleware(), handlers.ProfileUpdate)
// router.POST("/profile/change-password", utils.AuthMiddleware(), handlers.ChangePassword)

// router.GET("/profile/:id", utils.AuthRequired(), controllers.GetUserProfile)
// router.POST("/profile/change-password", utils.AuthRequired(), controllers.ChangePassword)
