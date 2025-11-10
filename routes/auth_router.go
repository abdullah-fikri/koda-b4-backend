package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	auth.POST("/register", controllers.RegisterUser)
	auth.POST("/login", controllers.LoginUser)
	auth.PUT("/update/:id", middleware.Auth(), controllers.UpdateUser)

	admin := r.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())
	admin.POST("/:id/picture", controllers.UploadPicture)
	admin.PUT("/:id/update", controllers.UpdateUser)
	admin.GET("/user", controllers.ListUser)
}
