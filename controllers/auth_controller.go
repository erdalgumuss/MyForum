package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id FROM users WHERE email = ?", input.Email).Scan(&user.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already taken"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	err = config.DB.QueryRow("SELECT id FROM users WHERE username = ?", input.Username).Scan(&user.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password.String), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	stmt, err := config.DB.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Username, input.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	sessionToken := uuid.New().String()
	c.SetCookie("session_token", sessionToken, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func ProcessLogin(c *gin.Context) {
	var input models.User
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid form data"})
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", input.Email).Scan(&user.ID, &user.Password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(input.Password)); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid email or password"})
		return
	}

	sessionToken := uuid.New().String()
	c.SetCookie("session_token", sessionToken, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/") // Redirect to the homepage
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func ProcessRegister(c *gin.Context) {
	var input models.User
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Invalid form data"})
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id FROM users WHERE email = ?", input.Email).Scan(&user.ID)
	if err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
		return
	}

	stmt, err := config.DB.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Username, input.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.Redirect(http.StatusFound, "/login") // Redirect to the login page if registration is successful
}

// ListUsers lists all users
func ListUsers(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, username, email, created_at, updated_at FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user row"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
