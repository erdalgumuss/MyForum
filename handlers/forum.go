package handlers

import (
	"MyForum/config"
	"MyForum/controllers"
	"MyForum/models"
	"log"
	"net/http"

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

	// Check Content-Type and handle accordingly
	if c.ContentType() == "application/x-www-form-urlencoded" {
		if err := c.ShouldBind(&input); err != nil {
			log.Println("Form data binding error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		log.Println("Form data binding successful:", input)

		// Call the controller function
		controllers.CreatePostWithPost(c, input)
	} else {
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println("JSON binding error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		log.Println("JSON binding successful:", input)

		// Call the controller function
		controllers.CreatePostWithPost(c, input)
	}
}

func GetPosts(c *gin.Context) {
	log.Println("GetPosts function called")

	rows, err := config.DB.Query("SELECT id, title, content, user_id, likes, dislikes FROM posts")
	if err != nil {
		log.Println("Failed to query posts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query posts"})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Likes, &post.Dislikes); err != nil {
			log.Println("Failed to scan post:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan post"})
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rows error"})
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
