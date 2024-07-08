package handlers

import (
	"log"
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func ProfileView(c *gin.Context) {
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		log.Println("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id, email, username FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		log.Println("Failed to retrieve user profile:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	log.Println("User profile retrieved for user ID:", user.ID)
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user": user,
	})
}

/*func ProfileUpdate(c *gin.Context) {
	userID := utils.GetUserIDFromSession(c)

	var input models.Profile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec("UPDATE users SET username = ?, email = ?, full_name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", input.Username, input.Email, input.FullName, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

/*func ChangePassword(c *gin.Context) {
	userID := utils.GetUserIDFromSession(c)

	var input models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedPassword string
	err := config.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&storedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(input.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	_, err = config.DB.Exec("UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", hashedPassword, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}*/
