package controllers

import (
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
	config.DB.Where("email = ?", input.Email).First(&user)
	if user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
		return
	}

	input.Password = string(hashedPassword)
	config.DB.Create(&input)

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	config.DB.Where("email = ?", input.Email).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
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

// ProcessLogin login formunu işler
func ProcessLogin(c *gin.Context) {
	// Burada login işlemleri yapılır
	// Örneğin, form verileri alınır ve doğrulama yapılır
	// Eğer başarılıysa kullanıcı ana sayfaya yönlendirilir
	// Başarısızsa tekrar login sayfasına yönlendirilir
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// ProcessRegister kayıt formunu işler
func ProcessRegister(c *gin.Context) {
	// Burada kayıt işlemleri yapılır
	// Örneğin, form verileri alınır ve yeni kullanıcı oluşturulur
	// Eğer başarılıysa kullanıcı login sayfasına yönlendirilir
	// Başarısızsa tekrar kayıt sayfasına yönlendirilir
}
