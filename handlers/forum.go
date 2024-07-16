package handlers

import (
	"log"
	"net/http"
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

	// Kullanıcı kimliğini context'ten al
	userID, ok := c.Get("userID")
	if !ok {
		log.Println("Kullanıcı kimliği post oturumda bulunamadı")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz"})
		return
	}

	// Content-Type'a göre veriyi bağla
	if c.ContentType() == "application/x-www-form-urlencoded" {
		if err := c.ShouldBind(&input); err != nil {
			log.Println("Form verisi bağlama hatası:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz giriş: " + err.Error()})
			return
		}
		log.Println("Form verisi başarılı şekilde bağlandı:", input)
	} else {
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println("JSON bağlama hatası:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz giriş: " + err.Error()})
			return
		}
		log.Println("JSON verisi başarılı şekilde bağlandı:", input)
	}

	// Posta kullanıcı kimliğini ve diğer bilgileri ekle
	input.UserID = userID.(int)
	input.CreatedAt = time.Now()

	// Controller fonksiyonunu çağır
	if err := controllers.CreatePostWithPost(input); err != nil {
		log.Println("Post oluşturulurken hata:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Post oluşturulamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post başarıyla oluşturuldu"})
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
