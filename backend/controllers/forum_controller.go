package controllers

import (
	"log"
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

func CreatePostWithPost(post models.Post) error {
	// Start a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Transaction start error:", err)
		return err
	}

	// Insert into posts table
	query := `
	INSERT INTO posts (user_id, title, content, image_url, created_at)
	VALUES (?, ?, ?, ?, ?)
	RETURNING id
	`
	var postID int
	err = tx.QueryRow(query, post.UserID, post.Title, post.Content, post.ImageURL, post.CreatedAt).Scan(&postID)
	if err != nil {
		log.Println("Error inserting post:", err)
		tx.Rollback()
		return err
	}

	// Insert categories into post_categories table
	categoryQuery := `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`
	for _, categoryID := range post.CategoryIDs {
		_, err := tx.Exec(categoryQuery, postID, categoryID)
		if err != nil {
			log.Println("Error inserting post category:", err)
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println("Transaction commit error:", err)
		return err
	}

	return nil
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
