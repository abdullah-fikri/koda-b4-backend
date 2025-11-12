package routes

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	AuthRoutes(r)
	ProductRouter(r)
	AdminRoutes(r)
	OrderRouter(r)
	FavoriteRouter(r)
	CartRoutes(r)
}
