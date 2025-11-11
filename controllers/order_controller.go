package controllers

import (
	"backend/lib"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Summary Create a new order
// @Description User membuat pesanan baru
// @Tags User - Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body models.CreateOrderRequest true "Order Request Body"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/order [post]
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

// OrderHistory godoc
// @Summary Get user's order history
// @Description Menampilkan list order milik user yang sedang login
// @Tags User - Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/history [get]
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

// OrderDetail godoc
// @Summary Get order detail by ID
// @Description User melihat detail 1 order miliknya
// @Tags User - Orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/order/{id} [get]
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
