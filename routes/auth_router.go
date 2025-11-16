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
    auth.POST("/forgot-password", controllers.ForgotPassword)
    auth.POST("/reset-password", controllers.ResetPassword)

	user := r.Group("/user")
	user.Use(middleware.Auth())
	user.GET("/profile", controllers.UserProfile)
	user.PUT("/profile/update", controllers.UpdateProfile)
	user.POST("/profile/upload", controllers.UploadUserPicture)
}
