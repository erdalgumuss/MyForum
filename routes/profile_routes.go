// routes/profile_routes.go

package routes

import (
	"MyForum/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterProfileRoutes profil rotalarını kaydeder
func ProfileRoutes(r *gin.Engine) {
	r.GET("/profile", controllers.GetProfile)
	r.PUT("/profile", controllers.UpdateProfile)
	r.POST("/profile/change-password", controllers.ChangePassword)
}
