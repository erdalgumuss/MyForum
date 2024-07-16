package routes

import (
	"MyForum/handlers"
	"MyForum/utils" // AuthMiddleware'i kullanmak için

	"github.com/gin-gonic/gin"
)

func ForumRoutes(r *gin.Engine) {
	// Oturum doğrulaması gerektiren rotalar
	authorized := r.Group("/")
	authorized.Use(utils.AuthMiddleware())
	{
		authorized.GET("/create-post", handlers.RenderCreatePostPage)
		authorized.POST("/create-post", handlers.CreatePost)
		authorized.POST("/comments", handlers.CreateComment)
		authorized.POST("/posts/:id/like", handlers.LikePost)
		authorized.POST("/posts/:id/dislike", handlers.DislikePost)
		authorized.POST("/comments/:id/like", handlers.LikeComment)
		authorized.POST("/comments/:id/dislike", handlers.DislikeComment)
	}

	// Oturum doğrulaması gerektirmeyen rotalar
	r.GET("/getpost", handlers.GetPosts)
	r.GET("/posts/:id", handlers.GetPost)
	r.GET("/forum", handlers.ShowForumPage)
	r.GET("/gallery", handlers.GalleryPage)
	r.GET("/rules", handlers.RulesPage)
	r.POST("/rules", handlers.RulesPage)
}
