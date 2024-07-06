package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// Register handles user registration.
func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email is already registered
	var existingUser models.User
	err := config.DB.QueryRow("SELECT id FROM users WHERE email = ?", input.Email).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already taken"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if username is already taken
	err = config.DB.QueryRow("SELECT id FROM users WHERE username = ?", input.Username).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user into database
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
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
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
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func ProcessRegister(c *gin.Context) {
	var input struct {
		Username string `form:"username" binding:"required"`
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	// Form verilerini al ve input yapısına bind et
	if err := c.ShouldBind(&input); err != nil {
		// Hatalı form verisi durumunda hata mesajı gönder
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Invalid form data"})
		return
	}

	// Veritabanında kullanıcı var mı diye kontrol et
	var user models.User
	err := config.DB.QueryRow("SELECT id FROM users WHERE email = ?", input.Email).Scan(&user.ID)
	if err == nil {
		// Eğer email zaten kayıtlıysa hata mesajı gönder
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Email already registered"})
		return
	} else if err != sql.ErrNoRows {
		// Veritabanı hatası durumunda hata mesajı gönder
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Database error"})
		fmt.Println("Database error:", err) // Hata konsola yazdırıldı
		return
	}

	// Şifreyi hash'le
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		// Şifreleme hatası durumunda hata mesajı gönder
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Password encryption failed"})
		fmt.Println("Password encryption failed:", err) // Hata konsola yazdırıldı
		return
	}

	// Kullanıcıyı veritabanına ekle
	stmt, err := config.DB.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		// Veritabanına hazırlık hatası durumunda hata mesajı gönder
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Database error"})
		fmt.Println("Prepare statement error:", err) // Hata konsola yazdırıldı
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Username, input.Email, hashedPassword)
	if err != nil {
		// Veritabanına ekleme hatası durumunda hata mesajı gönder
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Failed to create user"})
		fmt.Println("Exec statement error:", err) // Hata konsola yazdırıldı
		return
	}

	// Başarılı kayıt durumunda ana sayfaya yönlendir
	c.Redirect(http.StatusFound, "/")
}
