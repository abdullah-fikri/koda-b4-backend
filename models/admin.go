package models

import (
	"backend/config"
	"context"
	"time"
)

type OrderListItem struct {
	ID     int64     `json:"id"`
	Date   time.Time `json:"date"`
	Status string    `json:"status"`
	Total  float64   `json:"total"`
}

func GetAllOrders() ([]OrderListItem, error) {
	ctx := context.Background()

	query := `
	SELECT
	    o.id,
	    o.order_date,
	    o.status,
	    COALESCE(SUM(oi.subtotal), 0) AS total
	FROM orders o
	JOIN order_items oi ON oi.order_id = o.id
	GROUP BY o.id
	ORDER BY o.order_date DESC;
`

	rows, err := config.Db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []OrderListItem

	for rows.Next() {
		var o OrderListItem
		err = rows.Scan(&o.ID, &o.Date, &o.Status, &o.Total)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func UpdateOrderStatus(orderID int64, status string) error {
	ctx := context.Background()

	_, err := config.Db.Exec(ctx, `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`, status, orderID)

	if err != nil {
		return err
	}

	_, err = config.Db.Exec(ctx, `
		UPDATE order_items
		SET status = $1
		WHERE order_id = $2
	`, status, orderID)

	return err
}
