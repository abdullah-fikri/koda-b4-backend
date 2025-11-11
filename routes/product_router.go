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

	admin.GET("/product", controllers.AdminProductList)
	admin.POST("/product-create", controllers.CreateProduct)
	admin.PUT("/product/:id", controllers.UpdateProduct)
	admin.DELETE("/product/:id", controllers.DeleteProduct)
	admin.POST("/product/:id/pictures", controllers.UploadProductImages)

	r.GET("/products", controllers.Product)
	r.GET("/products/:id", controllers.ProductDetail)
}
