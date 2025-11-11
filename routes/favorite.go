package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func FavoriteRouter(r *gin.Engine){
	r.GET("/favorite-product", controllers.FavoriteProduct)
}