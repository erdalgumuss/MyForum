package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"MyForum/config"
	"MyForum/controllers"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"google.golang.org/api/idtoken"
)

var (
	facebookOauthConfig *oauth2.Config
	githubOauthConfig   *oauth2.Config
)

func ShowIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Home",
	})
}

func ProcessLogin(c *gin.Context) {
	controllers.Login(c)
}

// DM -- PM //
// ShowInbox renders the inbox.html template
func ShowInbox(c *gin.Context) {
	c.HTML(http.StatusOK, "inbox.html", gin.H{})
}

// DM -- PM //

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func ProcessRegister(c *gin.Context) {
	controllers.Register(c)
	//
}

func Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func GoogleLogin(c *gin.Context) {
	url := config.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Handle callback from Google OAuth after user approves access
func GoogleCallback(c *gin.Context) {
	// Exchange authorization code for OAuth token
	code := c.Query("code")
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Extract ID token from OAuth token
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("No id_token field in oauth2 token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token field in oauth2 token"})
		return
	}

	// Validate ID token to get user information
	payload, err := idtoken.Validate(context.Background(), idToken, config.GoogleOAuthConfig.ClientID)
	if err != nil {
		log.Println("Failed to validate id token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate id token"})
		return
	}

	// Extract user details from the validated payload
	email := payload.Claims["email"].(string)
	name := ""    // Initialize name as empty string
	surname := "" // Initialize surname as empty string

	// Extract name if available
	if val, ok := payload.Claims["name"].(string); ok {
		names := strings.Split(val, " ")
		if len(names) > 1 {
			name = names[0]
			surname = names[1]
		} else {
			name = val
		}
	}

	// Check if the user already exists in the database
	var user models.User
	err = config.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&user.ID)
	if err != nil {
		// If user doesn't exist, create a new user
		_, err = config.DB.Exec("INSERT INTO users (email, username, name, surname) VALUES (?, ?, ?, ?)", email, name+" "+surname, name, surname)
		if err != nil {
			log.Println("Failed to create user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Retrieve the newly created user's ID
		err = config.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&user.ID)
		if err != nil {
			log.Println("Failed to retrieve new user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve new user"})
			return
		}
	}

	// Create a session token for the user
	sessionToken, err := utils.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Set session token as a cookie
	c.SetCookie("session_token", sessionToken, 3600*24, "/", "localhost", false, true)

	// Redirect user to the desired page after successful login
	c.Redirect(http.StatusFound, "/")
}

func GitHubLogin(c *gin.Context) {
	url := githubOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := githubOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Println("Failed to get user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer response.Body.Close()

	var userInfo struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		log.Println("Failed to decode user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	var user models.User
	err = config.DB.QueryRow("SELECT id FROM users WHERE githubid = ?", userInfo.ID).Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = config.DB.Exec("INSERT INTO users (email, username, githubid) VALUES (?, ?, ?)", userInfo.Email, userInfo.Login, userInfo.ID)
			if err != nil {
				log.Println("Failed to create user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
			err = config.DB.QueryRow("SELECT id FROM users WHERE githubid = ?", userInfo.ID).Scan(&user.ID)
			if err != nil {
				log.Println("Failed to retrieve new user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve new user"})
				return
			}
		} else {
			log.Println("Failed to query user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	sessionToken, err := utils.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	c.SetCookie("session_token", sessionToken, 3600*24, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/")
}

func init() {
	config.LoadConfig()
	githubOauthConfig = &oauth2.Config{
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		RedirectURL:  config.GithubRedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	facebookOauthConfig = &oauth2.Config{
		ClientID:     config.FacebookClientID,
		ClientSecret: config.FacebookClientSecret,
		RedirectURL:  config.FacebookRedirectURL,
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
}

func FacebookLogin(c *gin.Context) {
	url := facebookOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func FacebookCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := facebookOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := facebookOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://graph.facebook.com/me?fields=id,name,email")
	if err != nil {
		log.Println("Failed to get user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer response.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		log.Println("Failed to decode user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	var user models.User
	err = config.DB.QueryRow("SELECT id FROM users WHERE facebookid = ?", userInfo.ID).Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = config.DB.Exec("INSERT INTO users (email, username, facebookid) VALUES (?, ?, ?)", userInfo.Email, userInfo.Name, userInfo.ID)
			if err != nil {
				log.Println("Failed to create user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
			err = config.DB.QueryRow("SELECT id FROM users WHERE facebookid = ?", userInfo.ID).Scan(&user.ID)
			if err != nil {
				log.Println("Failed to retrieve new user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve new user"})
				return
			}
		} else {
			log.Println("Failed to query user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	sessionToken, err := utils.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	c.SetCookie("session_token", sessionToken, 3600*24, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/")
}

//
