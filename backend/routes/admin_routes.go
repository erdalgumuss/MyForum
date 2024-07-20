package routes

import (
	"net/http"

	"MyForum/controllers"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	// Admin Routes
	admin := r.Group("/admin")
	{
		admin.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin.html", nil)
		})
		admin.POST("/login", controllers.ProcessAdminLogin)
		admin.Use(utils.AdminRequired()) // Admin oturum kontrolü
		admin.GET("/dashboard", controllers.AdminDashboard)
		admin.POST("/edit_user", controllers.EditUserRole)
		admin.POST("/assign_moderator", controllers.AssignModerator)
		admin.POST("/delete_user", controllers.DeleteUser)
		admin.GET("/user/:id", controllers.ViewUserProfile)
		admin.POST("/delete_post", controllers.DeletePost)
		admin.GET("/edit_post", controllers.EditPost)
		admin.POST("/update_post", controllers.UpdatePost)
		admin.GET("/request_moderator", controllers.ListModeratorRequests) // If needed separately
		admin.POST("/approve_moderator_request", controllers.ApproveModeratorRequest)
		admin.POST("/reject_moderator_request", controllers.RejectModeratorRequest)
	}
}
