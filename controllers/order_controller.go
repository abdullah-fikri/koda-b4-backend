package controllers

import (
	"backend/config"
	"backend/lib"
	"backend/models"
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Summary Create new order
// @Description Create new order from cart
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /orders [post]
func CreateOrder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int64)

	var req models.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return
	}

	//  wajib diisi
	if req.PaymentID == 0 || req.MethodID == 0 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "payment_id, shipping_id, and method_id are required",
		})
		return
	}

	if req.CustomerName == "" || req.CustomerPhone == "" || req.CustomerAddress == "" {
		var userData struct {
			Name    sql.NullString
			Phone   sql.NullString
			Address sql.NullString
		}

		err := config.Db.QueryRow(context.Background(), `
	SELECT username, phone, address 
	FROM profile 
	WHERE id = $1
`, userID).Scan(&userData.Name, &userData.Phone, &userData.Address)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Failed to fetch user info",
				Data:    err.Error(),
			})
			return
		}

		if !userData.Address.Valid || userData.Address.String == "" {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "complete the data first",
			})
			return
		}

		if req.CustomerName == "" && userData.Name.Valid {
			req.CustomerName = userData.Name.String
		}
		if req.CustomerPhone == "" && userData.Phone.Valid {
			req.CustomerPhone = userData.Phone.String
		}
		if req.CustomerAddress == "" && userData.Address.Valid {
			req.CustomerAddress = userData.Address.String
		}

	}

	order, err := models.CreateOrder(userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Failed to create order",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Order created successfully",
		Data:    order,
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

	month, _ := strconv.Atoi(ctx.DefaultQuery("month", "0"))
	shippingStr := ctx.DefaultQuery("shipping_id", "")
	if shippingStr == "" {
		shippingStr = ctx.DefaultQuery("shippings_id", "0")
	}
	shippingID, _ := strconv.Atoi(shippingStr)

	history, err := models.GetOrderHistoryByUserID(int64(user.Id), month, shippingID)
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
