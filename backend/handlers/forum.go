package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"MyForum/config"
	"MyForum/controllers"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func ShowForumPage(c *gin.Context) {
	c.HTML(http.StatusOK, "forum.html", gin.H{
		"Title": "Forum",
	})
}

func HandleForumPage(c *gin.Context) {
	// Handle form submission or other POST data processing here
	// For now, just render the forum page as a placeholder
	c.HTML(http.StatusOK, "forum.html", gin.H{
		"title":   "Forum",
		"message": "Post received",
	})
}

func RenderCreatePostPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_post.html", nil)
}

func GalleryPage(c *gin.Context) {
	c.HTML(http.StatusOK, "gallery.html", gin.H{
		"Title": "Galeri",
	})
}

func RulesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "rules.html", gin.H{
		"Title": "Kurallar",
	})
}

func CreatePost(c *gin.Context) {
	log.Println("CreatePost function called in handlers")

	var input models.Post

	// Get user ID from context
	userID, ok := c.Get("userID")
	if !ok {
		log.Println("Kullanıcı kimliği post oturumda bulunamadı")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz"})
		return
	}

	// Bind form data explicitly
	input.Title = c.PostForm("title")
	categoryNames := c.PostFormArray("categories") // Expecting an array of category names
	input.Content = c.PostForm("content")
	input.Username = c.PostForm("username")
	input.UserID = userID.(int)
	input.CreatedAt = time.Now()

	// Debug: log the received categories
	log.Printf("Received categories: %v\n", categoryNames)

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil {
		filename := filepath.Base(file.Filename)
		filepath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			log.Println("Dosya kaydedilirken hata:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya kaydedilemedi"})
			return
		}
		input.ImageURL = "/uploads/" + filename // Ensure this matches the URL path
		log.Printf("File uploaded successfully: %s\n", filepath)
	} else {
		log.Println("No file uploaded")
	}

	// Convert category names to IDs
	var categoryIDs []int
	for _, categoryName := range categoryNames {
		categoryName = strings.TrimSpace(categoryName)
		var categoryID int
		err := config.DB.QueryRow("SELECT id FROM categories WHERE name = ?", categoryName).Scan(&categoryID)
		if err != nil {
			log.Println("Invalid category name:", categoryName)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category name: " + categoryName})
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	input.CategoryIDs = categoryIDs

	log.Printf("Form data bound successfully: %+v\n", input)

	// Call controller function
	if err := controllers.CreatePostWithPost(input); err != nil {
		log.Println("Post oluşturulurken hata:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Post oluşturulamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post başarıyla oluşturuldu"})
}

func CreateComment(c *gin.Context) {
	// Get user ID from context
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse form data
	content := c.PostForm("content")
	postIDStr := c.PostForm("post_id")

	// Convert post_id to integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Fetch username from the database
	var username string
	err = config.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch username"})
		return
	}

	stmt, err := config.DB.Prepare("INSERT INTO comments (content, user_id, post_id, username, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, userID, postID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

func GetComments(c *gin.Context) {
	postID := c.Param("id")
	log.Println("Fetching comments for post ID:", postID)

	rows, err := config.DB.Query(`
		SELECT c.id, c.content, COALESCE(u.username, 'Unknown') AS username, c.post_id, c.likes, c.dislikes, c.created_at, c.updated_at 
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`, postID)
	if err != nil {
		log.Println("Failed to fetch comments from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	var comments []map[string]interface{}
	for rows.Next() {
		var comment models.Comment
		var username string
		var createdAt, updatedAt sql.NullTime
		err := rows.Scan(&comment.ID, &comment.Content, &username, &comment.PostID, &comment.Likes, &comment.Dislikes, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Failed to scan comment:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan comment"})
			return
		}
		comment.Username = username

		formattedCreatedAt := "Unknown"
		if createdAt.Valid {
			formattedCreatedAt = createdAt.Time.Format("January 2, 2006 at 3:04pm")
		}

		formattedUpdatedAt := "Unknown"
		if updatedAt.Valid {
			formattedUpdatedAt = updatedAt.Time.Format("January 2, 2006 at 3:04pm")
		}

		formattedComment := map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"username":   comment.Username,
			"post_id":    comment.PostID,
			"likes":      comment.Likes,
			"dislikes":   comment.Dislikes,
			"created_at": formattedCreatedAt,
			"updated_at": formattedUpdatedAt,
			"user_id":    comment.UserID,
		}
		comments = append(comments, formattedComment)
	}

	if err := rows.Err(); err != nil {
		log.Println("Failed to iterate comments:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate comments"})
		return
	}

	log.Println("Comments fetched:", comments)
	c.JSON(http.StatusOK, comments)
}

func getCategoryIDByName(categoryName string) (int, error) {
	var categoryID int
	query := "SELECT id FROM categories WHERE name = ?"

	err := config.DB.QueryRow(query, categoryName).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("category '%s' not found", categoryName)
		}
		log.Printf("Error fetching category ID for '%s': %v\n", categoryName, err)
		return 0, err
	}

	return categoryID, nil
}

func GetPosts(c *gin.Context) {
	var posts []models.Post

	// Query posts with sorting by created_at descending
	rows, err := config.DB.Query("SELECT id, COALESCE(username, '') AS username, title, content, user_id, likes, dislikes, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	// Iterate over rows and scan into Post structs
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Content, &post.UserID, &post.Likes, &post.Dislikes, &post.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan posts"})
			return
		}
		posts = append(posts, post)
	}

	// Handle any iteration errors
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through posts"})
		return
	}

	// Return posts as JSON response
	c.JSON(http.StatusOK, posts)
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

	// Query to fetch comments for the post
	rows, err := config.DB.Query(`
		SELECT c.id, c.content, COALESCE(u.username, 'Unknown') AS username, c.post_id, c.likes, c.dislikes, c.created_at, c.updated_at 
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`, id)
	if err != nil {
		log.Println("Failed to fetch comments from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	var comments []map[string]interface{}
	for rows.Next() {
		var comment models.Comment
		var username string
		var createdAt, updatedAt sql.NullTime
		err := rows.Scan(&comment.ID, &comment.Content, &username, &comment.PostID, &comment.Likes, &comment.Dislikes, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Failed to scan comment:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan comment"})
			return
		}
		comment.Username = username

		formattedCreatedAt := "Unknown"
		if createdAt.Valid {
			formattedCreatedAt = createdAt.Time.Format("January 2, 2006 at 3:04pm")
		}

		formattedUpdatedAt := "Unknown"
		if updatedAt.Valid {
			formattedUpdatedAt = updatedAt.Time.Format("January 2, 2006 at 3:04pm")
		}

		formattedComment := map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"username":   comment.Username,
			"post_id":    comment.PostID,
			"likes":      comment.Likes,
			"dislikes":   comment.Dislikes,
			"created_at": formattedCreatedAt,
			"updated_at": formattedUpdatedAt,
			"user_id":    comment.UserID,
		}
		comments = append(comments, formattedComment)
	}

	// Query to fetch category IDs and names
	rows, err = config.DB.Query(`
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
		"Comments":   comments,
	})
}

func DislikePost(c *gin.Context) {
}

func LikePost(c *gin.Context) {
}

func LikeComment(c *gin.Context) {
}

func DislikeComment(c *gin.Context) {
}
