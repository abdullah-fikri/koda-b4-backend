package controllers

import (
	"backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FavoriteProduct(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "4")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 4
	}

	products, total, err := models.Favorite(pageInt, limitInt)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "list favorite products",
		Data: models.Response{
			Success: true,
			Message: "list favorite products",
			Data: map[string]any{
				"products": products,
				"total":    total,
				"page":     pageInt,
				"limit":    limitInt,
			},
		},
	})
}
