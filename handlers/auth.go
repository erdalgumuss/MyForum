package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func ProcessLogin(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("JSON binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Received login request:", input)

	var storedUser models.User
	err := config.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", input.Email).Scan(&storedUser.ID, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No user found for email:", input.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			fmt.Println("Database query error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
		}
		return
	}

	fmt.Println("Stored user found:", storedUser)

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(input.Password))
	if err != nil {
		fmt.Println("Password comparison error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	sessionToken, err := utils.CreateSession(storedUser.ID)
	if err != nil {
		fmt.Println("Session creation error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.SetCookie("session_token", sessionToken, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	c.Redirect(http.StatusFound, "/") // Başarılı giriş sonrası yönlendirme
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func ProcessRegister(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO users (username, email, password, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", input.Username, input.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
