package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Register handles user registration.
func Register(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
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
	stmt, err := config.DB.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		fmt.Println("Failed to prepare statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Username, input.Email, hashedPassword)
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

/*func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}*/

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" `
		Password string `json:"password" `
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("ShouldBind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", input.Email).Scan(&user.ID, &user.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	sessionToken := uuid.New().String()
	c.SetCookie("session_token", sessionToken, 3600, "/profile", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
