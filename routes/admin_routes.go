package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middleware.Auth(), middleware.AdminOnly())

	//auth
	admin.POST("/user/:id/profile/upload", controllers.AdminUploadUserPicture)
	admin.PUT("/:id/update", controllers.AdminUpdateUser)
	admin.GET("/user", controllers.ListUser)

	//products
	admin.GET("/product", controllers.AdminProductList)
	admin.POST("/product-create", controllers.CreateProduct)
	admin.PUT("/product/:id", controllers.UpdateProduct)
	admin.DELETE("/product/:id", controllers.DeleteProduct)
	admin.POST("/product/:id/pictures", controllers.UploadProductImages)

	//order
	admin.GET("/orders", controllers.AdminOrderList)
	admin.PUT("/orders/:id/status", controllers.UpdateOrderStatus)
	admin.GET("/order/:id", controllers.OrderDetail)


	//category
	admin.POST("/category", controllers.CreateCategoryController)
	admin.GET("/category", controllers.GetAllCategoriesController)
	admin.PUT("/category/:id", controllers.UpdateCategoryController)
	admin.DELETE("/category/:id", controllers.DeleteCategoryController)
}
