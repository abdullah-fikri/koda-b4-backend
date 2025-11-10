package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRouter(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(middleware.Auth())

	user.POST("/order", controllers.CreateOrder)
	user.GET("/history", controllers.OrderHistory)
	user.GET("/order/:id", controllers.OrderDetail)
}
