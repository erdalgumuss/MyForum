package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"MyForum/config"
	"MyForum/controllers"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

func ShowForumPage(c *gin.Context) {
	c.HTML(http.StatusOK, "forum.html", gin.H{
		"Title": "Forum",
	})
}

// HandleForumPost handles POST requests to /forum
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
	log.Println("GetPosts function called")

	var posts []models.Post

	rows, err := config.DB.Query("SELECT id, COALESCE(username, '') AS username, title, content, user_id, likes, dislikes, created_at FROM posts")
	if err != nil {
		log.Println("Veritabanından postlar alınırken hata:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Postlar alınamadı"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Content, &post.UserID, &post.Likes, &post.Dislikes, &post.CreatedAt); err != nil {
			log.Println("Failed to scan post:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Postlar alınamadı"})
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Postlar alınamadı"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPostHandler handles the request and calls the controller function
func GetPost(c *gin.Context) {
	controllers.GetPost(c)
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
