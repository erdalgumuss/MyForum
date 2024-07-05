package routes

import (
	"MyForum/controllers"

	"github.com/gin-gonic/gin"
)

// AuthRoutes handles authentication related routes
func AuthRoutes(r *gin.Engine) {
	// GET requests
	r.GET("/login", controllers.ShowLoginPage)
	r.GET("/register", controllers.ShowRegisterPage)

	// POST requests
	r.POST("/login", controllers.ProcessLogin)
	r.POST("/register", controllers.Register)

	// Logout
	r.GET("/logout", controllers.Logout)
}
func Routes(r *gin.Engine) {
    // Auth routes
    r.GET("/login", controllers.ShowLoginPage)
    r.GET("/register", controllers.ShowRegisterPage)
    r.POST("/register", controllers.Register)
    r.POST("/login", controllers.ProcessLogin)
    r.GET("/logout", controllers.Logout)
    
    // List users route
    r.GET("/users", controllers.ListUsers)
}