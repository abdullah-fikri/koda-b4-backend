package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRouter(r *gin.Engine) {

	r.GET("/products", controllers.Product)
	r.GET("/products/:id", controllers.ProductDetail)
}
