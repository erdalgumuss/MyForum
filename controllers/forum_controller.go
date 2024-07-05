package controllers

import (
	"net/http"
	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var input models.Post
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
}

func CreateComment(c *gin.Context) {
	var input models.Comment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

func LikePost(c *gin.Context) {
	postID := c.Param("id")
	var post models.Post
	config.DB.First(&post, postID)
	post.Likes++
	config.DB.Save(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post liked"})
}

func DislikePost(c *gin.Context) {
	postID := c.Param("id")
	var post models.Post
	config.DB.First(&post, postID)
	post.Dislikes++
	config.DB.Save(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post disliked"})
}

func LikeComment(c *gin.Context) {
	commentID := c.Param("id")
	var comment models.Comment
	config.DB.First(&comment, commentID)
	comment.Likes++
	config.DB.Save(&comment)
	c.JSON(http.StatusOK, gin.H{"message": "Comment liked"})
}

func DislikeComment(c *gin.Context) {
	commentID := c.Param("id")
	var comment models.Comment
	config.DB.First(&comment, commentID)
	comment.Dislikes++
	config.DB.Save(&comment)
	c.JSON(http.StatusOK, gin.H{"message": "Comment disliked"})
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Preload("User").Preload("Categories").Find(&posts)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Posts": posts,
	})
}

func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := config.DB.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post": post,
	})
}
