package routes

import (
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.Use(middleware.CorsMiddleware())
	r.MaxMultipartMemory = 25 << 20
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"success": true,
			"message": "backend is running",
		})
	})
	AuthRoutes(r)
	ProductRouter(r)
	AdminRoutes(r)
	OrderRouter(r)
	FavoriteRouter(r)
	CartRoutes(r)
}
