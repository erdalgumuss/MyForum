package handlers

import (
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

// GetForums retrieves all forums
func GetForums(c *gin.Context) {
	var forums []models.Post
	config.DB.Preload("Categories").Find(&forums)
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

	forum := models.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID.(uint),
	}

	if len(input.Categories) > 0 {
		var categories []models.Category
		for _, categoryName := range input.Categories {
			var category models.Category
			if config.DB.Where("name = ?", categoryName).First(&category).RowsAffected == 0 {
				category = models.Category{Name: categoryName}
				config.DB.Create(&category)
			}
			categories = append(categories, category)
		}
		forum.Categories = categories
	}

	config.DB.Create(&forum)
	c.JSON(http.StatusOK, gin.H{"message": "Forum created successfully"})
}

// GetForum retrieves a specific forum by ID
func GetForum(c *gin.Context) {
	id := c.Param("id")

	var forum models.Post
	if config.DB.Preload("Categories").Preload("Comments").Where("id = ?", id).First(&forum).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forum not found"})
		return
	}

	c.HTML(http.StatusOK, "forum.html", gin.H{
		"forum": forum,
	})
}

// CreateComment handles comment creation
func CreateComment(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	comment := models.Comment{
		Content: input.Content,
		UserID:  userID.(uint),
		PostID:  utils.StringToUint(id), // Corrected here
	}

	config.DB.Create(&comment)
	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

// LikeForum handles liking a forum
func LikeForum(c *gin.Context) {
	id := c.Param("id")

	var forum models.Post
	if config.DB.Where("id = ?", id).First(&forum).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forum not found"})
		return
	}

	forum.Likes++
	config.DB.Save(&forum)

	c.JSON(http.StatusOK, gin.H{"message": "Forum liked successfully"})
}

// DislikeForum handles disliking a forum
func DislikeForum(c *gin.Context) {
	id := c.Param("id")

	var forum models.Post
	if config.DB.Where("id = ?", id).First(&forum).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forum not found"})
		return
	}

	forum.Dislikes++
	config.DB.Save(&forum)

	c.JSON(http.StatusOK, gin.H{"message": "Forum disliked successfully"})
}
