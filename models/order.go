package models

import (
	"backend/config"
	"context"
	"fmt"
	"time"
)

type CreateOrderRequest struct {
	PaymentID       int64  `json:"payment_id"`
	ShippingID      int64  `json:"shipping_id"`
	MethodID        int64  `json:"method_id"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerAddress string `json:"customer_address"`
}

type OrderDetail struct {
	ID              int64          `json:"order_id"`
	Invoice         string         `json:"invoice"`
	CustomerName    string         `json:"customer_name"`
	CustomerPhone   string         `json:"customer_phone"`
	CustomerAddress string         `json:"customer_address"`
	PaymentMethod   string         `json:"payment_method"`
	ShippingMethod  string         `json:"shipping_method"`
	OrderDate       time.Time      `json:"order_date"`
	Status          string         `json:"status"`
	Total           float64        `json:"total"`
	Items           []OrderItemRes `json:"items"`
}

type OrderItemRes struct {
	ProductName   string  `json:"product_name"`
	Variant       string  `json:"variant"`
	Size          string  `json:"size"`
	Qty           int     `json:"qty"`
	BasePrice     float64 `json:"base_price"`
	DiscountPrice float64 `json:"discount_price"`
	Img string `json:"img"`
}

type OrderResponse struct {
	OrderID         int64   `json:"order_id"`
	Invoice         string  `json:"invoice"`
	Total           float64 `json:"total"`
	CustomerName    string  `json:"customer_name"`
	CustomerPhone   string  `json:"customer_phone"`
	CustomerAddress string  `json:"customer_address"`
	Status          string  `json:"status"`
}

func CreateOrder(userID int64, req CreateOrderRequest) (OrderResponse, error) {
	ctx := context.Background()

	var cartItemCount int
	err := config.Db.QueryRow(ctx, `
		SELECT COUNT(ci.id)
		FROM cart_items ci
		JOIN cart c ON c.id = ci.cart_id
		WHERE c.user_id = $1
	`, userID).Scan(&cartItemCount)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to check cart: %w", err)
	}

	if cartItemCount == 0 {
		return OrderResponse{}, fmt.Errorf("cannot create order: cart is empty")
	}

	var total float64
	err = config.Db.QueryRow(ctx, `
	SELECT COALESCE(SUM(
		CASE 
			WHEN d.id IS NOT NULL 
				AND NOW() BETWEEN d.start_discount AND d.end_discount
			THEN (COALESCE(ps.price, p.price) - (COALESCE(ps.price, p.price) * d.percent_discount / 100)) * ci.qty
			ELSE COALESCE(ps.price, p.price) * ci.qty
		END
	), 0)
	FROM cart_items ci
	JOIN cart c ON c.id = ci.cart_id
	JOIN products p ON p.id = ci.product_id
	LEFT JOIN product_size ps ON ps.id = ci.size_id
	LEFT JOIN product_discount pd ON pd.product_id = ci.product_id
	LEFT JOIN discount d ON d.id = pd.discount_id
	WHERE c.user_id = $1
	`, userID).Scan(&total)
	if err != nil {
	return OrderResponse{}, fmt.Errorf("failed to calculate total: %w", err)
	}
	
	tx, err := config.Db.Begin(ctx)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	invoice := fmt.Sprintf("INV-%d-%d", time.Now().Unix(), userID)

	var orderID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (
			users_id, payment_id, shipping_id, method_id,
			order_date, customer_name, customer_phone, customer_address,
			total, invoice
		)
		VALUES ($1, $2, 3, $3, NOW(), $4, $5, $6, $7, $8)
		RETURNING id
	`,
		userID,
		req.PaymentID,
		req.MethodID,
		req.CustomerName,
		req.CustomerPhone,
		req.CustomerAddress,
		total,
		invoice,
	).Scan(&orderID)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to insert order: %w", err)
	}

	rows, err := tx.Query(ctx, `
	SELECT 
		ci.product_id,
		ci.variant_id,
		ci.size_id,
		ci.qty,
		COALESCE(ps.price, p.price) AS base_price,
		COALESCE(d.percent_discount, 0) AS discount_percent,
		CASE
			WHEN d.id IS NOT NULL
				AND NOW() BETWEEN d.start_discount AND d.end_discount
			THEN COALESCE(ps.price, p.price) - (COALESCE(ps.price, p.price) * COALESCE(d.percent_discount, 0) / 100)
			ELSE COALESCE(ps.price, p.price)
		END AS discount_price
	FROM cart_items ci
	JOIN cart c ON c.id = ci.cart_id
	JOIN products p ON p.id = ci.product_id
	LEFT JOIN product_size ps ON ps.id = ci.size_id
	LEFT JOIN product_discount pd ON pd.product_id = ci.product_id
	LEFT JOIN discount d ON d.id = pd.discount_id
	WHERE c.user_id = $1
	`, userID)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to fetch cart items: %w", err)
	}

	var items []struct {
		ProductID         int64
		VariantID, SizeID *int64
		Qty               int
		BasePrice         float64
		DiscountPercent   float64
		DiscountPrice     float64
	}

	for rows.Next() {
		var i struct {
			ProductID         int64
			VariantID, SizeID *int64
			Qty               int
			BasePrice         float64
			DiscountPercent   float64
			DiscountPrice     float64
		}
		if err := rows.Scan(&i.ProductID, &i.VariantID, &i.SizeID, &i.Qty, &i.BasePrice, &i.DiscountPercent, &i.DiscountPrice); err != nil {
			rows.Close()
			return OrderResponse{}, fmt.Errorf("failed to scan cart item: %w", err)
		}
		items = append(items, i)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return OrderResponse{}, fmt.Errorf("error while reading rows: %w", err)
	}

	for _, i := range items {
		subtotal := i.DiscountPrice * float64(i.Qty)
		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (
				order_id, product_id, variant_id, size_id, qty, 
				base_price, discount_price, discount_percent, subtotal
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, orderID, i.ProductID, i.VariantID, i.SizeID, i.Qty,
			i.BasePrice, i.DiscountPrice, i.DiscountPercent, subtotal)
		if err != nil {
			return OrderResponse{}, fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM cart_items 
		WHERE cart_id IN (SELECT id FROM cart WHERE user_id = $1)
	`, userID)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to clear cart: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return OrderResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return OrderResponse{
		OrderID:         orderID,
		Invoice:         invoice,
		Total:           total,
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		CustomerAddress: req.CustomerAddress,
		Status:          "On Progress",
	}, nil
}

func GetOrderHistoryByUserID(userID int64, month, shippingID, page, limit int) ([]map[string]interface{}, int, error) {
	ctx := context.Background()

	offset := (page - 1) * limit
	query := `
	SELECT 
		o.id AS order_id,
		o.invoice,
		o.order_date,
		COALESCE(o.total, 0) AS total,
		s.name AS shipping_status,
		COALESCE(MIN(pi.image), '') AS image
	FROM orders o
	JOIN shippings s ON s.id = o.shipping_id
	JOIN order_items oi ON oi.order_id = o.id
	JOIN products p ON p.id = oi.product_id
	LEFT JOIN product_img pi ON pi.product_id = p.id
	WHERE o.users_id = $1
	`

	args := []interface{}{userID}
	argIndex := 2

	if month > 0 {
		query += fmt.Sprintf(" AND EXTRACT(MONTH FROM o.order_date) = $%d", argIndex)
		args = append(args, month)
		argIndex++
	} else {
		query += " AND EXTRACT(MONTH FROM o.order_date) = (SELECT EXTRACT(MONTH FROM MAX(order_date)) FROM orders WHERE users_id = $1)"
	}

	if shippingID == 0 {
		shippingID = 3
	}
	query += fmt.Sprintf(" AND o.shipping_id = $%d", argIndex)
	args = append(args, shippingID)
	argIndex++

	query += `
	GROUP BY o.id, o.invoice, o.order_date, o.total, s.name
	ORDER BY o.order_date DESC
	LIMIT $%d OFFSET $%d
	`
	query = fmt.Sprintf(query, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := config.Db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var history []map[string]interface{}

	for rows.Next() {
		var (
			orderID   int64
			invoice   string
			orderDate time.Time
			total     float64
			status    string
			image     string
		)

		if err := rows.Scan(&orderID, &invoice, &orderDate, &total, &status, &image); err != nil {
			return nil, 0, err
		}

		history = append(history, map[string]interface{}{
			"order_id":   orderID,
			"invoice":    invoice,
			"order_date": orderDate,
			"total":      total,
			"status":     status,
			"image":      image,
		})
	}

	countQuery := `
		SELECT COUNT(DISTINCT o.id)
		FROM orders o
		WHERE o.users_id = $1
	`
	countArgs := []interface{}{userID}

	argIndex2 := 2
	if month > 0 {
		countQuery += fmt.Sprintf(" AND EXTRACT(MONTH FROM o.order_date) = $%d", argIndex2)
		countArgs = append(countArgs, month)
		argIndex2++
	} else {
		countQuery += " AND EXTRACT(MONTH FROM o.order_date) = (SELECT EXTRACT(MONTH FROM MAX(order_date)) FROM orders WHERE users_id = $1)"
	}

	countQuery += fmt.Sprintf(" AND o.shipping_id = $%d", argIndex2)
	countArgs = append(countArgs, shippingID)

	var totalItems int
	err = config.Db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	return history, totalItems, nil
}

func GetOrderDetail(orderID int64) (*OrderDetail, error) {
	ctx := context.Background()
	order := OrderDetail{}

	err := config.Db.QueryRow(ctx, `
		SELECT 
		o.id,
		o.invoice,
		o.customer_name,
		o.customer_phone,
		o.customer_address,
		o.order_date,
		s.name AS shipping_status,
		p.name AS payment_method,
		m.name AS shipping_method,
		o.total
	FROM orders AS o
	JOIN payment AS p ON o.payment_id = p.id
	JOIN shippings AS s ON o.shipping_id = s.id
	JOIN method AS m ON o.method_id = m.id    
	WHERE o.id = $1
`, orderID).Scan(
		&order.ID,
		&order.Invoice,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.CustomerAddress,
		&order.OrderDate,
		&order.Status,
		&order.PaymentMethod,
		&order.ShippingMethod,
		&order.Total,
	)

	if err != nil {
		return nil, err
	}

	rows, err := config.Db.Query(ctx, `
	SELECT 
		pr.name,
		COALESCE(v.name, '-') AS variant,
		COALESCE(sz.name, '-') AS size,
		oi.qty,
		oi.discount_price,
		COALESCE(pr.price, oi.discount_price) AS price,
		COALESCE(MIN(pimg.image), '') AS image
	FROM order_items oi
	JOIN products pr ON oi.product_id = pr.id
	LEFT JOIN variant v ON oi.variant_id = v.id
	LEFT JOIN size sz ON oi.size_id = sz.id
	LEFT JOIN product_img pimg ON pr.id = pimg.product_id
	WHERE oi.order_id = $1
	GROUP BY pr.name, v.name, sz.name, oi.qty, oi.discount_price, COALESCE(pr.price, oi.discount_price)
	`, orderID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item OrderItemRes
		if err := rows.Scan(&item.ProductName, &item.Variant, &item.Size, &item.Qty, &item.DiscountPrice,&item.BasePrice, &item.Img); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}
