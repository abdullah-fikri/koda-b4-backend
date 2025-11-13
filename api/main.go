package main

import (
	"backend/config"
	"backend/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	config.ConnectDb()
	config.Redis()

	app = gin.Default()
	routes.Routes(app)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
