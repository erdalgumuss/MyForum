package controllers

import (
	"log"
	"net/http"
	"time"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

func CreatePostWithPost(c *gin.Context, input models.Post) {
	log.Println("CreatePostWithPost function called in controllers")

	if config.DB == nil {
		log.Println("Database connection is nil in controller")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is nil"})
		return
	}
	log.Println("Database connection is OK in controller")

	stmt, err := config.DB.Prepare("INSERT INTO posts (title, content, user_id, username, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement in controller:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement: " + err.Error()})
		return
	}
	defer stmt.Close()
	log.Println("SQL statement prepared in controller")

	_, err = stmt.Exec(input.Title, input.Content, input.UserID, input.Username, time.Now())
	if err != nil {
		log.Println("Failed to execute statement in controller:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute statement: " + err.Error()})
		return
	}
	log.Println("SQL statement executed in controller")

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
	log.Println("Post created successfully in controller")
}

func CreatePost(c *gin.Context) {
	log.Println("CreatePost function called in controllers")

	var input models.Post
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("JSON binding error in controller:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	log.Println("JSON binding successful in controller:", input)

	if config.DB == nil {
		log.Println("Database connection is nil in controller")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is nil"})
		return
	}
	log.Println("Database connection is OK in controller")

	stmt, err := config.DB.Prepare("INSERT INTO posts (title, content, user_id, username) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement in controller:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement: " + err.Error()})
		return
	}
	defer stmt.Close()
	log.Println("SQL statement prepared in controller")

	_, err = stmt.Exec(input.Title, input.Content, input.UserID, input.Username)
	if err != nil {
		log.Println("Failed to execute statement in controller:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute statement: " + err.Error()})
		return
	}
	log.Println("SQL statement executed in controller")

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
	log.Println("Post created successfully in controller")
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

// GetPost retrieves a single post by ID
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	err := config.DB.QueryRow("SELECT id, title, content, likes, dislikes, user_id, username, created_at FROM posts WHERE id = ?", id).
		Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.UserID, &post.Username, &post.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	log.Printf("Retrieved Post: %+v\n", post) // Log the post data

	c.HTML(http.StatusOK, "post.html", gin.H{"Post": post})
}
