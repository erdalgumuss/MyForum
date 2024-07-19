package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/*
INSERT INTO users (username, email, password, name, surname, created_at, updated_at, githubid, role) VALUES (

	'admin',
	'admin@example.com',
	'efaa948ea3425eca1978671c8c1b9d2d',
	'Admin',
	'User',
	DATETIME('now'),
	DATETIME('now'),
	NULL,
	'admin'

);
*/
type User struct {
	ID       int
	Username string
	Email    string
	Role     string
}

func AdminDashboard(c *gin.Context) {
	// Fetch user data
	users, err := fetchUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Fetch moderator requests
	requests, err := fetchModeratorRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch moderator requests"})
		return
	}

	// Assuming you have the username and role stored in session or context
	username := "admin_username" // replace with actual username from session/context
	role := "admin_role"         // replace with actual role from session/context

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"username": username,
		"role":     role,
		"users":    users,
		"requests": requests,
	})
}

func DeleteUser(c *gin.Context) {
	userID := c.PostForm("user_id")
	_, err := config.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ViewUserProfile(c *gin.Context) {
	userID := c.Param("id")

	var user User
	err := config.DB.QueryRow("SELECT id, username, email, role FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	rows, err := config.DB.Query("SELECT id, title, content FROM posts WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		posts = append(posts, post)
	}

	c.HTML(http.StatusOK, "user_profile.html", gin.H{
		"user":  user,
		"posts": posts,
	})
}

func ProcessAdminLogin(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user struct {
		ID       int
		Username string
		Password string
		Role     string
	}

	query := "SELECT id, username, password, role FROM users WHERE username = ? AND role = 'admin'"
	err := config.DB.QueryRow(query, input.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	// Admin successfully authenticated
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	session.Set("role", user.Role)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func ListUsers(c *gin.Context) {
	var users []models.User

	rows, err := config.DB.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning users"})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func EditUser(c *gin.Context) {
	var input struct {
		UserID int    `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
		// Add other editable fields here
	}

	// Bind JSON input to struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the user from database
	var user models.User
	query := "SELECT id, username, email, role FROM users WHERE id = ?"
	err := config.DB.QueryRow(query, input.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Update the user's role
	user.Role = input.Role

	// Execute update query
	updateQuery := "UPDATE users SET role = ? WHERE id = ?"
	_, err = config.DB.Exec(updateQuery, user.Role, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// admin_controller.go

func AssignModerator(c *gin.Context) {
	var input struct {
		UserID int `json:"user_id" binding:"required"`
	}

	// Bind JSON input to struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the user from database
	var user models.User
	query := "SELECT id, username, email, role FROM users WHERE id = ?"
	err := config.DB.QueryRow(query, input.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Update the user's role to "moderator"
	user.Role = "moderator"

	// Execute update query
	updateQuery := "UPDATE users SET role = ? WHERE id = ?"
	_, err = config.DB.Exec(updateQuery, user.Role, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "User assigned as moderator"})
}

func EditUserRole(c *gin.Context) {
	var input struct {
		UserID int    `form:"user_id" binding:"required"`
		Role   string `form:"role" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Kullanıcı rolünü güncellemek için gerekli veritabanı işlemleri burada yapılmalı
	// Örneğin:
	// query := "UPDATE users SET role = ? WHERE id = ?"
	// _, err := config.DB.Exec(query, input.Role, input.UserID)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ListModeratorRequests(c *gin.Context) {
	// Adding a log to see if the function is being called
	fmt.Println("ListModeratorRequests function called")

	rows, err := config.DB.Query("SELECT id, user_id, status, created_at FROM moderator_requests WHERE status = 'pending'")
	if err != nil {
		fmt.Println("Error querying database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}
	defer rows.Close()

	var requests []struct {
		ID        int
		UserID    int
		Status    string
		CreatedAt string
	}

	for rows.Next() {
		var req struct {
			ID        int
			UserID    int
			Status    string
			CreatedAt string
		}
		err := rows.Scan(&req.ID, &req.UserID, &req.Status, &req.CreatedAt)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse requests"})
			return
		}
		requests = append(requests, req)
	}

	// Adding a log to see the fetched requests
	fmt.Println("Fetched requests:", requests)

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{"requests": requests})
}

func ApproveModeratorRequest(c *gin.Context) {
	var input struct {
		RequestID int `form:"request_id" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	var userID int
	err := config.DB.QueryRow("SELECT user_id FROM moderator_requests WHERE id = ?", input.RequestID).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found: " + err.Error()})
		return
	}

	_, err = config.DB.Exec("UPDATE users SET role = 'moderator' WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role: " + err.Error()})
		return
	}

	_, err = config.DB.Exec("UPDATE moderator_requests SET status = 'approved' WHERE id = ?", input.RequestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request approved"})
}

func RejectModeratorRequest(c *gin.Context) {
	var input struct {
		RequestID int `json:"request_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	_, err := config.DB.Exec("UPDATE moderator_requests SET status = 'rejected' WHERE id = ?", input.RequestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request rejected"})
}

func fetchUsers() ([]struct {
	ID       int
	Username string
	Email    string
	Role     string
}, error,
) {
	rows, err := config.DB.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []struct {
		ID       int
		Username string
		Email    string
		Role     string
	}

	for rows.Next() {
		var user struct {
			ID       int
			Username string
			Email    string
			Role     string
		}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func fetchModeratorRequests() ([]struct {
	ID        int
	UserID    int
	Status    string
	CreatedAt string
}, error,
) {
	rows, err := config.DB.Query("SELECT id, user_id, status, created_at FROM moderator_requests WHERE status = 'pending'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []struct {
		ID        int
		UserID    int
		Status    string
		CreatedAt string
	}

	for rows.Next() {
		var req struct {
			ID        int
			UserID    int
			Status    string
			CreatedAt string
		}
		err := rows.Scan(&req.ID, &req.UserID, &req.Status, &req.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}
