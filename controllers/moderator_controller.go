package controllers

import (
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ModeratorDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "moderator_dashboard.html", gin.H{
		"title": "Moderator Dashboard",
	})
}

func PendingPosts(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, title, content FROM posts WHERE status = 'pending'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		posts = append(posts, post)
	}

	c.HTML(http.StatusOK, "pending_posts.html", gin.H{
		"posts": posts,
	})
}

func ApprovePost(c *gin.Context) {
	postID := c.PostForm("post_id")
	var content string
	err := config.DB.QueryRow("SELECT content FROM posts WHERE id = ?", postID).Scan(&content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if utils.ContainsForbiddenContent(content) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post contains forbidden content"})
		return
	}

	_, err = config.DB.Exec("UPDATE posts SET status = 'approved' WHERE id = ?", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func RejectPost(c *gin.Context) {
	postID := c.PostForm("post_id")
	_, err := config.DB.Exec("UPDATE posts SET status = 'rejected' WHERE id = ?", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeletePost(c *gin.Context) {
	postID := c.PostForm("post_id")
	_, err := config.DB.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func RequestModerator(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	_, err := config.DB.Exec("INSERT INTO moderator_requests (user_id, status) VALUES (?, 'pending')", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func EditPost(c *gin.Context) {
	postID := c.Query("post_id")

	var post models.Post
	err := config.DB.QueryRow("SELECT id, title, content FROM posts WHERE id = ?", postID).Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	c.HTML(http.StatusOK, "edit_post.html", gin.H{
		"post": post,
	})
}

func UpdatePost(c *gin.Context) {
	postID := c.PostForm("post_id")
	title := c.PostForm("title")
	content := c.PostForm("content")

	_, err := config.DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", title, content, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
