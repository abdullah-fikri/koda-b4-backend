package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRouter(r *gin.Engine) {

	admin := r.Group("/admin")
	admin.Use(middleware.Auth()) 
	admin.Use(middleware.AdminOnly())    

	admin.POST("/product", controllers.CreateProduct)
	admin.PUT("/product/:id", controllers.UpdateProduct)

	r.GET("/products", controllers.Product)
	r.GET("/products/:id", controllers.ProductDetail)
}
