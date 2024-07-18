package controllers

import (
	"log"
	"net/http"
	"strings"

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

func GetPost(c *gin.Context) {
	// Extract post ID from URL parameter
	id := c.Param("id")

	var post models.Post
	var categoryIDs []int
	var categoryNames []string

	// Query to fetch post details including username from users table
	err := config.DB.QueryRow(`
		SELECT p.id, p.title, p.content, p.likes, p.dislikes, p.user_id, p.image_url, u.username
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.UserID, &post.ImageURL, &post.Username)

	if err != nil {
		log.Println("Error fetching post:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Query to fetch category IDs and names
	rows, err := config.DB.Query(`
		SELECT c.id, c.name
		FROM categories c
		INNER JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
	`, post.ID)
	if err != nil {
		log.Println("Error fetching categories:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categories"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var categoryID int
		var categoryName string
		if err := rows.Scan(&categoryID, &categoryName); err != nil {
			log.Println("Error scanning category:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning categories"})
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
		categoryNames = append(categoryNames, categoryName)
	}

	post.CategoryIDs = categoryIDs

	// Convert category names slice to a comma-separated string for display purposes
	categoriesString := strings.Join(categoryNames, ", ")

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post":       post,
		"Categories": categoriesString,
	})
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
