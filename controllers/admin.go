package controllers

import (
	"backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func UpdateOrderStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: "Invalid order ID"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}

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
