package routes

import (
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.Use(middleware.CorsMiddleware())
	r.MaxMultipartMemory = 25 << 20
	AuthRoutes(r)
	ProductRouter(r)
	AdminRoutes(r)
	OrderRouter(r)
	FavoriteRouter(r)
	CartRoutes(r)
}
