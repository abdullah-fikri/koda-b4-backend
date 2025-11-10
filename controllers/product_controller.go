package controllers

import (
	"backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Product(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")
	sort := ctx.DefaultQuery("sort", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	products, total, err := models.GetProducts(page, limit, search, sort)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	totalPage := int((total + int64(limit) - 1) / int64(limit))

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "products fetched",
		Data: map[string]any{
			"products": products,
			"pagination": map[string]any{
				"page":       page,
				"limit":      limit,
				"total_data": total,
				"total_page": totalPage,
			},
		},
	})
}

func ProductDetail(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product id",
		})
		return
	}

	product, err := models.GetProductByID(productID)
	if err != nil {
		c.JSON(400,
			models.Response{
				Success: false,
				Message: "product not found",
			})
		return
	}

	c.JSON(200, models.Response{
		Success: true,
		Message: "success",
		Data: map[any]any{
			"data": product,
		},
	})
}

func CreateProduct(ctx *gin.Context) {
	var req models.CreateProductRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	product, err := models.CreateProduct(req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Product created",
		Data:    product,
	})
}

func UpdateProduct(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	var req models.CreateProductRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: err.Error()})
		return
	}

	product, err := models.UpdateProduct(int64(id), req)
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product updated",
		Data:    product,
	})
}
