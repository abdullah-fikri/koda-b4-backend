package models

import (
	"backend/config"
	"context"
)

type ReqCart struct {
	UserID    int64   `json:"user_id"`
	ProductID int64   `json:"product_id"`
	VariantID *int64  `json:"variant_id,omitempty"`
	SizeID    *int64  `json:"size_id,omitempty"`
	Qty       int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type CartItem struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	VariantID *int64  `json:"variant_id,omitempty"`
	SizeID    *int64  `json:"size_id,omitempty"`
	Qty       int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type CartItemResponse struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"name"`
	Variant     string  `json:"variant"`
	Size        string  `json:"size"`
	Price       float64 `json:"price"`
	Qty         int     `json:"quantity"`
	Image       string  `json:"image"`
}
func AddToCart(req ReqCart) (int64, error) {
	ctx := context.Background()
	var cartID int64

	err := config.Db.QueryRow(ctx, `
		SELECT id FROM cart WHERE user_id = $1
	`, req.UserID).Scan(&cartID)
	if err != nil {
		err = config.Db.QueryRow(ctx, `
			INSERT INTO cart (user_id) VALUES ($1) RETURNING id
		`, req.UserID).Scan(&cartID)
		if err != nil {
			return 0, err
		}
	}

	// cek
	var existingID int64
	checkErr := config.Db.QueryRow(ctx, `
		SELECT id FROM cart_items 
		WHERE cart_id = $1
		  AND product_id = $2
		  AND COALESCE(variant_id, 0) = COALESCE($3::BIGINT, 0)
		  AND COALESCE(size_id, 0) = COALESCE($4::BIGINT, 0)
	`, cartID, req.ProductID, req.VariantID, req.SizeID).Scan(&existingID)

	var price float64
	if req.SizeID != nil {
		err = config.Db.QueryRow(ctx, `
			SELECT price FROM product_size WHERE id = $1
		`, req.SizeID).Scan(&price)
	} else {
		err = config.Db.QueryRow(ctx, `
			SELECT MIN(price) FROM product_size WHERE product_id = $1
		`, req.ProductID).Scan(&price)
	}
	if err != nil {
		return 0, err
	}

	// update 
	if checkErr == nil {
		_, err = config.Db.Exec(ctx, `
			UPDATE cart_items
			SET qty = qty + $1::INT,
			    subtotal = (qty + $1::INT) * $2::NUMERIC
			WHERE id = $3::BIGINT
		`, req.Qty, price, existingID)
		if err != nil {
			return 0, err
		}
		return cartID, nil
	}

	//  buat baru
	_, err = config.Db.Exec(ctx, `
		INSERT INTO cart_items (cart_id, product_id, variant_id, size_id, qty, subtotal)
		VALUES ($1::BIGINT, $2::BIGINT, $3::BIGINT, $4::BIGINT, $5::INT, $6::NUMERIC * $5::NUMERIC)
	`, cartID, req.ProductID, req.VariantID, req.SizeID, req.Qty, price)
	if err != nil {
		return 0, err
	}

	return cartID, nil
}
