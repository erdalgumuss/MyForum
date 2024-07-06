package handlers

import (
	"net/http"

	"MyForum/controllers"

	"github.com/gin-gonic/gin"
)

func ShowIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func ProcessLogin(c *gin.Context) {
	controllers.ProcessLogin(c)
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func ProcessRegister(c *gin.Context) {
	controllers.ProcessRegister(c)
}

func Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
