// routes/profile_routes.go

package routes

import (
	"MyForum/controllers"

	"github.com/gin-gonic/gin"
)

// ProfileRoutes profil rotalarını kaydeder
func ProfileRoutes(r *gin.Engine) {
	r.GET("/profile", controllers.GetUserProfile)
	r.PUT("/profile", controllers.UpdateProfile)
	r.POST("/profile/change-password", controllers.ChangePassword)
}
