package controllers

import (
	"backend/models"

	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user, err := models.Register(req)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Register success",
		Data:    user,
	})
}
