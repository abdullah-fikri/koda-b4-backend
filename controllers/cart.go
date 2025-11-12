package controllers

import (
	"backend/lib"
	"backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToCart(ctx *gin.Context) {
	userData, _ := ctx.Get("user")
	user := userData.(lib.UserPayload)
	userID := int64(user.Id)

	var req models.ReqCart
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if req.Qty <= 0 {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Quantity must be greater than 0",
		})
		return
	}

	req.UserID = userID

	cartID, err := models.AddToCart(req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Item added to cart",
		Data:    cartID,
	})
}

func GetCart(ctx *gin.Context) {
	userData, _ := ctx.Get("user")
	user := userData.(lib.UserPayload)
	userID := int64(user.Id)

	carts, err := models.GetCart(userID)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Cart successfully",
		Data:    carts,
	})
}
func DeleteCart(ctx *gin.Context) {
	userData, _ := ctx.Get("user")
	user := userData.(lib.UserPayload)
	userID := int64(user.Id)

	cartItemIDStr := ctx.Param("id")
	cartItemID, err := strconv.ParseInt(cartItemIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid cart item id",
		})
		return
	}

	err = models.DeleteCartItem(userID, cartItemID)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product removed from cart",
	})
}
