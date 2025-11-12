package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func CartRoutes(r *gin.Engine){
	cart:= r.Group("/cart")
	cart.Use(middleware.Auth())

	cart.POST("/", controllers.AddToCart)
	cart.GET("/", controllers.GetCart)
}