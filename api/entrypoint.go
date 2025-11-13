package handler

import (
	"backend/config"
	"backend/models"
	"backend/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var App *gin.Engine

func init() {
	App = gin.New()
	App.Use(gin.Recovery())

	router := App.Group("/")
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Backend is running well",
		})
	})

	routes.Routes(App)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if config.Db == nil {
		os.Getenv("DATABASE_URL")
		config.ConnectDb()
	}
	App.ServeHTTP(w, r)
}
