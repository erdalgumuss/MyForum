package routes

import (
	"MyForum/controllers"
	"github.com/gin-gonic/gin"
)


func ForumRoutes(r *gin.Engine) {
	r.GET("/", controllers.GetPosts)
	r.GET("/posts/:id", controllers.GetPost)
	r.POST("/posts", controllers.CreatePost)
	r.POST("/comments", controllers.CreateComment)
	r.POST("/posts/:id/like", controllers.LikePost)
	r.POST("/posts/:id/dislike", controllers.DislikePost)
	r.POST("/comments/:id/like", controllers.LikeComment)
	r.POST("/comments/:id/dislike", controllers.DislikeComment)
}
