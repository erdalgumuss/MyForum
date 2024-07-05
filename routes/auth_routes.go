package routes

import (

"time"
   "MyForum/controllers"
    "github.com/gin-gonic/gin"
)
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func AuthRoutes(r *gin.Engine) {
    // Rotaları tanımla
    r.GET("/login", controllers.ShowLoginPage)
    r.POST("/login", controllers.ProcessLogin)
    r.GET("/register", controllers.ShowRegisterPage)
    r.POST("/register", controllers.ProcessRegister)
    r.POST("/logout", controllers.Logout)
}
