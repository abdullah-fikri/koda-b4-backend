package api

import (
	"backend/config"
	"backend/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

var App *gin.Engine

func init() {
	config.ConnectDb()
	config.Redis()

	App = gin.Default()

	routes.Routes(App)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	App.ServeHTTP(w, r)
}
