package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	// Oturum açmış kullanıcıyı al
	userID, ok := c.Get("userID")
	if !ok {
		log.Println("Kullanıcı kimliği oturumda bulunamadı")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz"})
		return
	}

	/// Kullanıcı bilgilerini veritabanından al
	var user models.User
	err := config.DB.QueryRow("SELECT id, email, username, name, surname FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.Username, &user.Name, &user.Surname)
	if err != nil {
		log.Println("Kullanıcı profilini getirme başarısız oldu:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı profilini getirme başarısız oldu"})
		return
	}

	log.Println("Kullanıcı profili getirildi, Kullanıcı ID:", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"name":     user.Name,
		"surname":  user.Surname,
	})
}

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

func GetUserPosts(c *gin.Context) {
	userID := c.Param("id")
	var posts []models.Post

	// Query to fetch posts and their category IDs
	rows, err := config.DB.Query(`
		SELECT p.id, p.title, p.content, p.user_id, p.image_url, p.likes, p.dislikes, 
			GROUP_CONCAT(pc.category_id) AS category_ids
		FROM posts p
		LEFT JOIN post_categories pc ON p.id = pc.post_id
		WHERE p.user_id = ?
		GROUP BY p.id, p.title, p.content, p.user_id, p.image_url, p.likes, p.dislikes
	`, userID)
	if err != nil {
		log.Println("Error fetching user posts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user posts"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var categoryIDs string
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.ImageURL, &post.Likes, &post.Dislikes, &categoryIDs); err != nil {
			log.Println("Error scanning post row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user posts"})
			return
		}

		// Convert the comma-separated category IDs string to a slice of integers
		if categoryIDs != "" {
			for _, idStr := range strings.Split(categoryIDs, ",") {
				id, err := strconv.Atoi(idStr)
				if err == nil {
					post.CategoryIDs = append(post.CategoryIDs, id)
				}
			}
		}

		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

func GetUserComments(c *gin.Context) {
	userID := c.Param("id")
	var comments []struct {
		models.Comment
		PostTitle string `json:"post_title"`
	}

	rows, err := config.DB.Query(`
        SELECT c.id, c.content, c.post_id, c.user_id, c.username, c.likes, c.dislikes, c.created_at, c.updated_at, p.title 
        FROM comments c
        JOIN posts p ON c.post_id = p.id
        WHERE c.user_id = ?
    `, userID)
	if err != nil {
		log.Println("Failed to fetch user comments from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user comments"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment struct {
			models.Comment
			PostTitle string `json:"post_title"`
		}
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.Username, &comment.Likes, &comment.Dislikes, &comment.CreatedAt, &comment.UpdatedAt, &comment.PostTitle); err != nil {
			log.Println("Failed to scan user comments:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user comments"})
			return
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating through user comments:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through user comments"})
		return
	}

	log.Println("User comments fetched:", comments)
	c.JSON(http.StatusOK, comments)
}

func GetUserLikes(c *gin.Context) {
	userID := c.Param("id")
	var likes []models.Like

	rows, err := config.DB.Query(`
		SELECT ul.id, ul.user_id, ul.post_id, ul.comment_id, ul.action, p.title
		FROM user_likes ul
		LEFT JOIN posts p ON ul.post_id = p.id
		WHERE ul.user_id = ?
	`, userID)
	if err != nil {
		log.Println("Failed to fetch user likes from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user likes"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var like models.Like
		var postID sql.NullInt64
		var postTitle sql.NullString
		if err := rows.Scan(&like.ID, &like.UserID, &postID, &like.CommentID, &like.Action, &postTitle); err != nil {
			log.Println("Failed to scan user likes:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user likes"})
			return
		}
		like.PostID = int(postID.Int64)
		like.PostTitle = postTitle.String
		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating through user likes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through user likes"})
		return
	}

	log.Println("User likes fetched:", likes)
	c.JSON(http.StatusOK, likes)
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
}*/

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
