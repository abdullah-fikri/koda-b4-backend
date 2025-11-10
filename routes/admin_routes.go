package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())

	admin.GET("/orders", controllers.AdminOrderList)
	admin.PUT("/orders/:id/status", controllers.UpdateOrderStatus)

}
