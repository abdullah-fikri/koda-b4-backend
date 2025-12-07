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

type CartItemResponse struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"name"`
	Variant     string  `json:"variant"`
	Size        string  `json:"size"`
	Price       float64 `json:"base-price"`
	Subtotal    float64 `json:"subtotal"`
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

	// update
	if checkErr == nil {
		_, err = config.Db.Exec(ctx, `
			UPDATE cart_items
			SET qty = qty + $1::INT
			WHERE id = $2::BIGINT
		`, req.Qty, existingID)
		if err != nil {
			return 0, err
		}
		return cartID, nil
	}

	// buat baru
	_, err = config.Db.Exec(ctx, `
		INSERT INTO cart_items (cart_id, product_id, variant_id, size_id, qty)
		VALUES ($1::BIGINT, $2::BIGINT, $3::BIGINT, $4::BIGINT, $5::INT)
	`, cartID, req.ProductID, req.VariantID, req.SizeID, req.Qty)
	if err != nil {
		return 0, err
	}

	return cartID, nil
}

func GetCart(userID int64) ([]CartItemResponse, error) {
	ctx := context.Background()

	var cartID int64
	err := config.Db.QueryRow(ctx, `
		SELECT id FROM cart WHERE user_id = $1
	`, userID).Scan(&cartID)
	if err != nil {
		return []CartItemResponse{}, nil
	}

	rows, err := config.Db.Query(ctx, `
	SELECT 
		ci.id AS cart_item_id,
		p.id AS product_id,
		p.name AS product_name,
		COALESCE(v.name, '') AS variant_name,
		COALESCE(s.name, '') AS size_name,
		COALESCE(ps.price, p.price) AS price,
		ci.qty,
		(
			SELECT pi.image 
			FROM product_img pi 
			WHERE pi.product_id = p.id 
			LIMIT 1
		) AS image
	FROM cart_items ci
	JOIN products p ON p.id = ci.product_id
	LEFT JOIN variant v ON v.id = ci.variant_id
	LEFT JOIN size s ON s.id = ci.size_id
	LEFT JOIN product_size ps ON ps.product_id = p.id AND ps.size_id = ci.size_id
	WHERE ci.cart_id = $1
	`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []CartItemResponse

	for rows.Next() {
		var item CartItemResponse

		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.Variant,
			&item.Size,
			&item.Price,
			&item.Qty,
			&item.Image,
		)
		if err != nil {
			return nil, err
		}
		item.Subtotal = item.Price * float64(item.Qty)

		carts = append(carts, item)
	}

	return carts, nil
}

func DeleteCartItem(userID, cartItemID int64) error {
	ctx := context.Background()

	query := `
		DELETE FROM cart_items 
		WHERE id = $1 
		AND cart_id = (
			SELECT id FROM cart WHERE user_id = $2
		)
	`
	_, err := config.Db.Exec(ctx, query, cartItemID, userID)
	return err
}
