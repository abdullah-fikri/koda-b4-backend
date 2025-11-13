package main

import (
	"backend/config"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs"
)

// @title CoffeeShop API
// @version 1.0
// @description API documentation for Coffee Shop App
// @host localhost:8082
// @BasePath /
func main() {
	godotenv.Load()
	config.ConnectDb()
	config.Redis()
	r := gin.Default()

	routes.Routes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8082")
}
