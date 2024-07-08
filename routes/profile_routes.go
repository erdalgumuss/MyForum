package routes

import (
	"MyForum/handlers"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(router *gin.Engine) {
	router.GET("/profile", utils.AuthMiddleware(), handlers.ShowIndexPage)
	router.PUT("/profile", utils.AuthMiddleware(), handlers.ProfileUpdate)
	router.POST("/profile/change-password", utils.AuthMiddleware(), handlers.ChangePassword)
}
