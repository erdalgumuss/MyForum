package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/*INSERT INTO users (username, email, password, name, surname, created_at, updated_at, githubid, role) VALUES (
	'admin',
	'admin@example.com',
	'$2a$10$4bvg9T55V370.Z5mKhc3jeN54.OgnGG9pnjJ6r3y98Cbj02bfAKdm',
	'Admin',
	'User',
	DATETIME('now'),
	DATETIME('now'),
	NULL,
	'admin'
);
*/
// Register handles user registration.
func Register(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"}) ////git
		fmt.Println("Invalid JSON provided:", err)
		return
	}

	fmt.Println("Received registration data:", input)

	// Check if email is already registered
	var existingUser models.User
	err := config.DB.QueryRow("SELECT id FROM users WHERE email = ?", input.Email).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already taken"})
		fmt.Println("Email already taken:", input.Email)
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		fmt.Println("Database error:", err)
		return
	}

	// Check if username is already taken
	err = config.DB.QueryRow("SELECT id FROM users WHERE username = ?", input.Username).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		fmt.Println("Username already taken:", input.Username)
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		fmt.Println("Database error:", err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		fmt.Println("Failed to hash password:", err)
		return
	}

	// Insert user into database
	stmt, err := config.DB.Prepare("INSERT INTO users(username, name, surname, email, password) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		fmt.Println("Failed to prepare statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Username, input.Name, input.Surname, input.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		fmt.Println("Failed to create user:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Logout handles user logout.
func Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("ShouldBind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserByEmail(input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Mail veya Şifre Yanlış"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Invalidate old sessions
	err = utils.InvalidateOldSessions(user.ID)
	if err != nil {
		log.Println("Error invalidating old sessions:", err)
	}

	// Create a new session
	sessionToken, err := utils.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	c.SetCookie("session_token", sessionToken, 3600, "/", "localhost", false, true)

	// Set user session with role
	utils.SetUserSession(c, *user) // Dereference the pointer here

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"surname": user.Surname,
			"email":   user.Email,
		},
	})
}

// GoogleLogin handles the Google login redirection.

func getUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.QueryRow("SELECT id, email, username, name, surname, role, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Username, &user.Name, &user.Surname, &user.Role, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
