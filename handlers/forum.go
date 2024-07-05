package handlers

import (
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

// GetForums retrieves all forums
func GetForums(c *gin.Context) {
	rows, err := config.DB.Query("SELECT * FROM posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var forums []models.Post
	for rows.Next() {
		var forum models.Post
		err := rows.Scan(&forum.ID, &forum.Title, &forum.Content, &forum.UserID, &forum.Likes, &forum.Dislikes, &forum.CreatedAt, &forum.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		forums = append(forums, forum)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"forums": forums,
	})
}

// CreateForum handles forum creation
func CreateForum(c *gin.Context) {
	var input struct {
		Title      string   `json:"title" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		Categories []string `json:"categories"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	// Örnek olarak, INSERT işlemi için SQL sorgusu
	result, err := config.DB.Exec("INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)", input.Title, input.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	forumID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Forum created successfully", "forum_id": forumID})
}
