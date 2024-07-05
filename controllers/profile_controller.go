package controllers

import (
	"net/http"

	"MyForum/config"
	"MyForum/models"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	user.Email = input.Email
	user.Username = input.Username

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
// GetProfile kullanıcı profilini gösterir
func GetProfile(c *gin.Context) {
    // Kullanıcıyı veritabanından al ve profil sayfasını göster
    c.HTML(http.StatusOK, "profile.html", gin.H{})
}
// UpdateProfile kullanıcı profilini günceller
func UpdateProfile(c *gin.Context) {
    // Kullanıcı profilini güncelleyen işlemler yapılır
    // Örneğin, form verileri alınır ve veritabanında güncelleme yapılır
    c.JSON(http.StatusOK, gin.H{
        "message": "Profil güncellendi",
    })
}
// ChangePassword kullanıcı şifresini değiştirir
func ChangePassword(c *gin.Context) {
    // Kullanıcı şifresini değiştiren işlemler yapılır
    // Örneğin, form verileri alınır ve veritabanında şifre güncelleme işlemi yapılır
    c.JSON(http.StatusOK, gin.H{
        "message": "Şifre değiştirildi",
    })
}
