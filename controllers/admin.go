package controllers

import (
	"backend/models"
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
	orders, err := models.GetAllOrders()
	if err != nil {
		ctx.JSON(500, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Data:    orders,
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
