package controllers

import (
	"database/sql"
	"log"
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
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, exists := c.Get("role")
	if !exists || role == "moderator" || role == "admin" {
		log.Printf("Invalid role: %v, or role not found in context", role)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already a moderator or admin"})
		return
	}

	log.Printf("User ID: %v is requesting to be a moderator", userID)

	_, err := config.DB.Exec("INSERT INTO moderator_requests (user_id, status) VALUES (?, 'pending')", userID)
	if err != nil {
		log.Println("Failed to submit moderator request:", err)
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
func AssignModerator(c *gin.Context) {
	var input struct {
		UserID int `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	query := "SELECT id, username, email, role FROM users WHERE id = ?"
	err := config.DB.QueryRow(query, input.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	user.Role = sql.NullString{String: "moderator", Valid: true}

	updateQuery := "UPDATE users SET role = ? WHERE id = ?"
	_, err = config.DB.Exec(updateQuery, user.Role.String, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Update the role in the session
	session := sessions.Default(c)
	session.Set("role", user.Role.String)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User assigned as moderator"})
}
