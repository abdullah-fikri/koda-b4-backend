package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"encoding/json"
	"strconv"
	"time"

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
	key := ctx.Request.RequestURI
	var cacheData struct {
		Products   []models.Product `json:"products"`
		Pagination map[string]any   `json:"pagination"`
	}

	cache, err := config.Rdb.Get(context.Background(), key).Result()
	if err == nil && cache != "" {
		_ = json.Unmarshal([]byte(cache), &cacheData)
		ctx.JSON(200, models.Response{
			Success: true,
			Message: "list favorite products ( from cache )",
			Data:    cacheData,
		})
		return
	}

	products, total, err := models.Favorite(pageInt, limitInt)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	totalPage := int((total + int64(limitInt) - 1) / int64(limitInt))

	cacheData = struct {
		Products   []models.Product `json:"products"`
		Pagination map[string]any   `json:"pagination"`
	}{
		Products: products,
		Pagination: map[string]any{
			"page":       page,
			"limit":      limit,
			"total_data": total,
			"total_page": totalPage,
		},
	}

	jsonData, err := json.Marshal(cacheData)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "failed marshal data",
		})
		return
	}
	config.Rdb.Set(context.Background(), key, jsonData, 15*time.Minute)

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "list products favorite",
		Data:    cacheData,
	})
}
