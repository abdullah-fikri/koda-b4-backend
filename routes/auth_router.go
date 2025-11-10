package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/auth/register", controllers.RegisterUser)
	r.POST("/auth/login", controllers.LoginUser)
}
