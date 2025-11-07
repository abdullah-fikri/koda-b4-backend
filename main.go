package main

import (
	"backend/config"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDb()
	r := gin.Default()
	routes.Routes(r)

	r.Run(":8081")
}
