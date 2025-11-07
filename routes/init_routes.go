package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.POST("/auth/register", controllers.RegisterUser)
}
