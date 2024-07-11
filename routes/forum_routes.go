package routes

import (
	"MyForum/handlers" // controllers yerine handlers

	"github.com/gin-gonic/gin"
)

func ForumRoutes(r *gin.Engine) {
	r.GET("/", handlers.ShowForumPage)
	r.GET("/getpost", handlers.GetPosts)  // controllers yerine handlers
	r.GET("/posts/:id", handlers.GetPost) // controllers yerine handlers
	r.POST("/posts", handlers.CreatePost) // controllers yerine handlers
	r.POST("/create-post", handlers.CreatePost)
	r.GET("/create-post", handlers.RenderCreatePostPage)
	r.POST("/comments", handlers.CreateComment)              // controllers yerine handlers
	r.POST("/posts/:id/like", handlers.LikePost)             // controllers yerine handlers
	r.POST("/posts/:id/dislike", handlers.DislikePost)       // controllers yerine handlers
	r.POST("/comments/:id/like", handlers.LikeComment)       // controllers yerine handlers
	r.POST("/comments/:id/dislike", handlers.DislikeComment) // controllers yerine handlers
}
