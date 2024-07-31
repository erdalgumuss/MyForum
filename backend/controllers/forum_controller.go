package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

func CreatePostWithPost(post models.Post) (int, error) {
	// Start a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Transaction start error:", err)
		return 0, err
	}

	// Insert into posts table
	query := `
	INSERT INTO posts (user_id, title, content, username, image_url, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	RETURNING id
	`
	var postID int
	err = tx.QueryRow(query, post.UserID, post.Title, post.Content, post.Username, post.ImageURL, post.CreatedAt.Format("2006-01-02 15:04:05")).Scan(&postID)
	if err != nil {
		log.Println("Error inserting post:", err)
		tx.Rollback()
		return 0, err
	}

	// Insert categories into post_categories table
	categoryQuery := `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`
	for _, categoryID := range post.CategoryIDs {
		_, err := tx.Exec(categoryQuery, postID, categoryID)
		if err != nil {
			log.Println("Error inserting post category:", err)
			tx.Rollback()
			return 0, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println("Transaction commit error:", err)
		return 0, err
	}

	return postID, nil
}

func EditPostHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	query := "SELECT id, title, content FROM posts WHERE id = ?"
	err = config.DB.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		}
		return
	}

	c.HTML(http.StatusOK, "edit_post.html", gin.H{"post": post})
}

func UpdatePostHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.PostForm("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")

	query := "UPDATE posts SET title = ?, content = ? WHERE id = ?"
	_, err = config.DB.Exec(query, title, content, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
