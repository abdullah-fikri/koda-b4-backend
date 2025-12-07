package controllers

import (
	"backend/lib"
	"backend/models"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminOrderList godoc
// @Summary Get all orders (Admin)
// @Description Admin melihat semua order
// @Tags Admin - Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/orders [get]
func AdminOrderList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	orders, totalItems, err := models.GetAllOrders(page, limit)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	totalPage := int((totalItems + int64(limit) - 1) / int64(limit))

	rawQuery := ctx.Request.URL.Query()
	extraQuery := url.Values{}
	for k, v := range rawQuery {
		extraQuery[k] = append([]string{}, v...)
	}

	baseURL := os.Getenv("APP_BASE_URL")
	path := ctx.Request.URL.Path

	links := lib.Hateoas(baseURL, path, page, limit, totalPage, extraQuery)
	pagination := lib.Pagination(page, limit, totalPage, totalItems, links)

	ctx.JSON(200, models.Response{
		Success:    true,
		Message:    "list all order",
		Pagination: pagination,
		Data:       orders,
	})
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Admin mengubah status pesanan
// @Tags Admin - Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "Update Status"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/orders/{id}/status [put]
func UpdateOrderStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: "Invalid order ID"})
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: err.Error()})
		return
	}

	err = models.UpdateOrderStatus(int64(orderID), req.Status)
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Order status updated",
	})
}


// create category
func CreateCategoryController(ctx *gin.Context){
	var c models.Categories
	if err := ctx.ShouldBindJSON(&c); err !=nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	category, err := models.CreateCategory(c.Name)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Category Created",
		Data: category,
	})
}

// get all categories
func GetAllCategoriesController(c *gin.Context) {
	data, err := models.GetAllCategories()
	if err != nil {
		c.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, models.Response{
		Success: true,
		Message: "list all categories",
		Data:    data,
	})
}


// update category
func UpdateCategoryController(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var body models.Categories

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: err.Error()})
		return
	}

	updated, err := models.UpdateCategory(id, body.Name)
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Category updated",
		Data:    updated,
	})
}


// delete 
func DeleteCategoryController(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	err := models.DeleteCategory(id)
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Category deleted",
	})
}
