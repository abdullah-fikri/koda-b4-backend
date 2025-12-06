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

type Categories struct {
	Id int `json:id`
	Name string `json:"name" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status int `json:"status"`
}
func GetAllOrders() ([]OrderListItem, error) {
	ctx := context.Background()

	query := `
		SELECT
		    o.id,
		    o.order_date,
		    s.name AS shipping_name,
		    COALESCE(SUM(oi.subtotal), 0) AS total
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		LEFT JOIN shippings s ON s.id = o.shipping_id
		GROUP BY o.id, s.name
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

func UpdateOrderStatus(orderID int64, status int) error {
	ctx := context.Background()

	_, err := config.Db.Exec(ctx, `
		UPDATE orders
		SET shipping_id = $1
		WHERE id = $2
	`, status, orderID)

	if err != nil {
		return err
	}

	return err
}


func CreateCategory(name string) (Categories, error) {
    ctx := context.Background()

    var c Categories

    err := config.Db.QueryRow(ctx,
        "INSERT INTO categories (name) VALUES ($1) RETURNING name",
        name,
    ).Scan(&c.Name)

    if err != nil {
        return Categories{}, err
    }

    return c, nil
}

func GetAllCategories() ([]Categories, error) {
	ctx := context.Background()
	query := "SELECT id,name FROM categories ORDER BY id DESC"

	rows, err := config.Db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []Categories
	for rows.Next(){
		var c Categories
		err = rows.Scan(&c.Id,&c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories,nil
}


func UpdateCategory(id int, name string) (Categories, error) {
	ctx := context.Background()

	query := `
		UPDATE categories 
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name
	`

	var c Categories

	err := config.Db.QueryRow(ctx, query, name, id).Scan(&c.Id, &c.Name)
	if err != nil {
		return Categories{}, err
	}

	return c, nil
}

func DeleteCategory(id int) error {
	ctx := context.Background()

	query := `DELETE FROM categories WHERE id=$1`

	_, err := config.Db.Exec(ctx, query, id)
	return err
}
