package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

// Google OAuth 2.0 configuration
var googleOauthConfig = &oauth2.Config{
	ClientID:     "YOUR_GOOGLE_CLIENT_ID",
	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Endpoint:     google.Endpoint,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
}

// Register handles user registration.
func Register(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})//
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

// Login handles user login.
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

	sessionToken, err := utils.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	c.SetCookie("session_token", sessionToken, 3600, "/", "localhost", false, true)
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
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the callback from Google OAuth 2.0.
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("No id_token field in oauth2 token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token field in oauth2 token"})
		return
	}

	payload, err := idtoken.Validate(context.Background(), idToken, googleOauthConfig.ClientID)
	if err != nil {
		log.Println("Failed to validate id token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate id token"})
		return
	}

	userID := payload.Subject
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	var user models.User
	err = config.DB.QueryRow("SELECT id FROM users WHERE google_id = ?", userID).Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = config.DB.Exec("INSERT INTO users (google_id, email, username) VALUES (?, ?, ?)", userID, email, name)
			if err != nil {
				log.Println("Failed to create user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			err = config.DB.QueryRow("SELECT id FROM users WHERE google_id = ?", userID).Scan(&user.ID)
			if err != nil {
				log.Println("Failed to retrieve new user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		} else {
			log.Println("Failed to query user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	sessionToken, err := utils.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	c.SetCookie("session_token", sessionToken, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func getUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Password)
	return &user, err
}
