package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middleware.Auth())
	admin.Use(middleware.AdminOnly())

	admin.POST("/register", controllers.RegisterAd)
	admin.POST("/update", controllers.UpdateUserAd)

	r.POST("/auth/register", controllers.RegisterUser)
	r.POST("/auth/login", controllers.LoginUser)
	r.PUT("/auth/update", controllers.UpdateUser)
}
