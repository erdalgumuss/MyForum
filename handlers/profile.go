package handlers

import (
	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProfileView handles viewing user profile
func ProfileView(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user": user,
	})
}

// ProfileUpdate handles updating user profile
func ProfileUpdate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	user.Email = input.Email
	user.Username = input.Username
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func RegisterProfileRoutes(router *gin.Engine) {
	router.GET("/profile", utils.AuthRequired(), ProfileView)
	router.POST("/profile", utils.AuthRequired(), ProfileUpdate)
}
