package handlers

import (
	"context"
	"log"
	"net/http"

	"MyForum/config"
	"MyForum/controllers"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

func ShowIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Home",
	})
}

func ProcessLogin(c *gin.Context) {
	controllers.Login(c)
}

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

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
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

	payload, err := idtoken.Validate(context.Background(), idToken, config.GoogleOAuthConfig.ClientID)
	if err != nil {
		log.Println("Failed to validate id token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate id token"})
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	var user models.User
	err = config.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&user.ID)
	if err != nil {
		_, err = config.DB.Exec("INSERT INTO users (email, username) VALUES (?, ?)", email, name)
		if err != nil {
			log.Println("Failed to create user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		err = config.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&user.ID)
		if err != nil {
			log.Println("Failed to retrieve new user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve new user"})
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
