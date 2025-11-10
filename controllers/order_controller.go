package controllers

import (
	"backend/lib"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateOrder(ctx *gin.Context) {
	userData, _ := ctx.Get("user")
	user := userData.(lib.UserPayload)

	userID := int64(user.Id)

	var req models.CreateOrderRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error()})
		return
	}

	orderID, err := models.CreateOrder(userID, req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Order created",
		Data:    orderID,
	})
}

func OrderHistory(ctx *gin.Context) {
	userData, _ := ctx.Get("user")
	user := userData.(lib.UserPayload)

	history, err := models.GetOrderHistoryByUserID(int64(user.Id))
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Data:    history,
	})
}

func OrderDetail(ctx *gin.Context) {
	idParam := ctx.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid order ID"})
		return
	}

	result, err := models.GetOrderDetail(int64(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error()})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Data:    result})
}
