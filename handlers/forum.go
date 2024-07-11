package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowForumPage(c *gin.Context) {
	c.HTML(http.StatusOK, "/forum.html", gin.H{
		"Title": "Forum",
	})
}

func GetPosts(c *gin.Context) {
}

func GetPost(c *gin.Context) {
}

func CreatePost(c *gin.Context) {
}

func DislikePost(c *gin.Context) {
}

func CreateComment(c *gin.Context) {
}

func LikePost(c *gin.Context) {
}

func LikeComment(c *gin.Context) {
}

func DislikeComment(c *gin.Context) {
}
