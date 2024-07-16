package controllers

import (
	"log"
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

func CreatePostWithPost(post models.Post) error {
	query := `
	INSERT INTO posts (user_id, title, content, image_url, created_at)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := config.DB.Exec(query, post.UserID, post.Title, post.Content, post.ImageURL, post.CreatedAt)
	if err != nil {
		log.Println("Veritabanına post kaydedilirken hata:", err)
		return err
	}

	return nil
}

// CreateComment handles the creation of a new comment
func CreateComment(c *gin.Context) {
	var input models.Comment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := config.DB.Prepare("INSERT INTO comments (content, user_id, post_id) VALUES (?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Content, input.UserID, input.PostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

// LikePost handles the liking of a post
func LikePost(c *gin.Context) {
	postID := c.Param("id")
	result, err := config.DB.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked"})
}

// DislikePost handles the disliking of a post
func DislikePost(c *gin.Context) {
	postID := c.Param("id")
	result, err := config.DB.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike post"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post disliked"})
}

// LikeComment handles the liking of a comment
func LikeComment(c *gin.Context) {
	commentID := c.Param("id")
	result, err := config.DB.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like comment"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment liked"})
}

// DislikeComment handles the disliking of a comment
func DislikeComment(c *gin.Context) {
	commentID := c.Param("id")
	result, err := config.DB.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?", commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike comment"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment disliked"})
}

// GetPosts retrieves all posts
func GetPosts(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, title, content, likes, dislikes, user_id, username FROM posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.UserID, &post.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan post"})
			return
		}
		posts = append(posts, post)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"Posts": posts})
}

func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	err := config.DB.QueryRow("SELECT id, title, content, likes, dislikes, user_id, image_url FROM posts WHERE id = ?", id).
		Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.UserID, &post.ImageURL)
	if err != nil {
		log.Println("Error fetching post:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.HTML(http.StatusOK, "post.html", gin.H{"Post": post})
}
