package models

import (
	"backend/config"
	"context"
	"time"
)

type OrderItemReq struct {
	ProductID int64   `json:"product_id"`
	VariantID *int64  `json:"variant_id,omitempty"`
	SizeID    *int64  `json:"size_id,omitempty"`
	Qty       int     `json:"qty"`
	Subtotal  float64 `json:"subtotal"`
}

type CreateOrderRequest struct {
	PaymentID       int64          `json:"payment_id"`
	ShippingID      int64          `json:"shipping_id"`
	CustomerName    string         `json:"customer_name"`
	CustomerPhone   string         `json:"customer_phone"`
	CustomerAddress string         `json:"customer_address"`
	Items           []OrderItemReq `json:"items"`
}

type OrderDetail struct {
	ID              int64          `json:"order_id"`
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
	ProductName string  `json:"product_name"`
	Variant     string  `json:"variant"`
	Size        string  `json:"size"`
	Qty         int     `json:"qty"`
	Subtotal    float64 `json:"subtotal"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"Done"`
}


func CreateOrder(userID int64, req CreateOrderRequest) (int64, error) {
	ctx := context.Background()

	var orderID int64
	err := config.Db.QueryRow(ctx,
		`INSERT INTO orders (users_id, payment_id, shipping_id, customer_name, customer_phone, customer_address, order_date, status)
	 VALUES ($1, $2, $3, $4, $5, $6, now(), 'pending')
	 RETURNING id`,
		userID, req.PaymentID, req.ShippingID, req.CustomerName, req.CustomerPhone, req.CustomerAddress,
	).Scan(&orderID)

	if err != nil {
		return 0, err
	}

	for _, item := range req.Items {
		_, err = config.Db.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, variant_id, size_id, qty, subtotal, status)
			 VALUES ($1,$2,$3,$4,$5,$6,'pending')`,
			orderID, item.ProductID, item.VariantID, item.SizeID, item.Qty, item.Subtotal,
		)
		if err != nil {
			return 0, err
		}
	}

	return orderID, nil
}
func GetOrderHistoryByUserID(userID int64) ([]map[string]interface{}, error) {
	ctx := context.Background()

	rows, err := config.Db.Query(ctx,
		`SELECT 
            o.id, o.order_date, o.status,
            oi.product_id, oi.qty, oi.subtotal
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        WHERE o.users_id = $1
        ORDER BY o.order_date DESC`,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []map[string]interface{}

	for rows.Next() {
		var (
			orderID   int64
			orderDate time.Time
			status    string
			productID int64
			qty       int
			subtotal  float64
		)

		if err := rows.Scan(&orderID, &orderDate, &status, &productID, &qty, &subtotal); err != nil {
			return nil, err
		}

		history = append(history, map[string]interface{}{
			"order_id":   orderID,
			"order_date": orderDate,
			"status":     status,
			"product_id": productID,
			"qty":        qty,
			"subtotal":   subtotal,
		})
	}

	return history, nil
}

func GetOrderDetail(orderID int64) (*OrderDetail, error) {
	ctx := context.Background()
	order := OrderDetail{}

	err := config.Db.QueryRow(ctx, `
    SELECT 
        o.id,
        o.customer_name,
        o.customer_phone,
        o.customer_address,
        o.order_date,
        o.status,
        p.name,
        s.name,
        COALESCE(SUM(oi.subtotal), 0) AS total
    FROM orders AS o
    JOIN payment AS p ON o.payment_id = p.id
    JOIN shippings AS s ON o.shipping_id = s.id
    JOIN order_items oi ON oi.order_id = o.id
    WHERE o.id = $1
    GROUP BY o.id, p.name, s.name
`, orderID).Scan(
		&order.ID,
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
			oi.subtotal
		FROM order_items oi
		JOIN products pr ON oi.product_id = pr.id
		LEFT JOIN variant v ON oi.variant_id = v.id
		LEFT JOIN size sz ON oi.size_id = sz.id
		WHERE oi.order_id = $1
	`, orderID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item OrderItemRes
		if err := rows.Scan(&item.ProductName, &item.Variant, &item.Size, &item.Qty, &item.Subtotal); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}
