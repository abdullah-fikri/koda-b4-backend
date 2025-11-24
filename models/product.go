package models

import (
	"backend/config"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ProductSize struct {
	SizeID   int64   `json:"size_id" binding:"required,gt=0"`
	SizeName string  `json:"size_name"`
	Price    float64 `json:"price" binding:"required,gte=0"`
}

type ProductAdmin struct {
	ID          int64  `json:"id"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Sizes       string `json:"sizes"`
	Method      string `json:"method"`
	Stock       int64  `json:"stock"`
}

type Product struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	MinPrice    float64       `json:"min_price"`
	Stock       int64         `json:"stock"`
	Category    string        `json:"category"`
	Images      []string      `json:"images"`
	Sizes       []ProductSize `json:"sizes"`
	Variants    []string      `json:"variants"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string        `json:"name" binding:"required,min=1,max=200"`
	Description string        `json:"description"`
	Stock       int           `json:"stock" binding:"required,gte=0"`
	CategoryID  int           `json:"category_id"`
	Images      []string      `json:"images"`
	Variants    []int         `json:"variants"`
	Sizes       []ProductSize `json:"sizes"`
}

// admin version
func GetProductsAdmin(page, limit int, search string) ([]ProductAdmin, int64, error) {
	ctx := context.Background()

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}

	offset := (page - 1) * limit

	query := `
	SELECT
    p.id,
    COALESCE(pi.image, '') AS image,
    p.name,
    p.description,
    COALESCE(MIN(ps.price), 0) AS price,
    COALESCE(string_agg(DISTINCT s.name, ', '), '') AS sizes,
    COALESCE(string_agg(DISTINCT m.name, ', '), '') AS methods,
    p.stock
FROM products p
LEFT JOIN (
    SELECT DISTINCT ON (product_id) product_id, image
    FROM product_img
    ORDER BY product_id, id ASC
) pi ON pi.product_id = p.id
LEFT JOIN product_size ps ON ps.product_id = p.id
LEFT JOIN size s ON s.id = ps.size_id
LEFT JOIN product_method pm ON pm.product_id = p.id
LEFT JOIN method m ON m.id = pm.method_id
WHERE p.name ILIKE '%' || $1 || '%'
GROUP BY p.id, pi.image, p.name, p.description, p.stock
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3`



	rows, err := config.Db.Query(ctx, query, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []ProductAdmin

	for rows.Next() {
		var p ProductAdmin
		err := rows.Scan(&p.ID, &p.Image, &p.Name, &p.Description, &p.Price, &p.Sizes, &p.Method, &p.Stock)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}

	var total int64
	err = config.Db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE name ILIKE '%' || $1 || '%'`, search).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// user version
func GetProducts(page, limit int, search, sort string, minPrice int, maxPrice int, categoryIDs []int) ([]Product, int64, error) {
	ctx := context.Background()

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	orderByClause := "p.created_at DESC"
	switch sort {
	case "oldest":
		orderByClause = "p.created_at ASC"
	case "price_low":
		orderByClause = "min_price ASC"
	case "price_high":
		orderByClause = "min_price DESC"
	case "name_asc":
		orderByClause = "p.name ASC"
	case "name_desc":
		orderByClause = "p.name DESC"
	}

	query := `
SELECT
	p.id,
	p.name,
	p.description,
	COALESCE(
		(SELECT MIN(price) FROM product_size WHERE product_id = p.id),
		p.price
	  ) AS min_price,
	p.stock,
	COALESCE(c.name, '') AS category,
	COALESCE(json_agg(DISTINCT pi.image) FILTER (WHERE pi.image IS NOT NULL), '[]') AS images,
	COALESCE(json_agg(DISTINCT v.name) FILTER (WHERE v.name IS NOT NULL), '[]') AS variants,
	COALESCE(
		json_agg(DISTINCT jsonb_build_object(
			'size_id', s.id,
			'size_name', s.name,
			'price', ps.price
		)) FILTER (WHERE s.id IS NOT NULL),
		'[]'
	) AS sizes,
	p.created_at,
	p.updated_at
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN product_img pi ON pi.product_id = p.id
LEFT JOIN product_variant pv ON pv.product_id = p.id
LEFT JOIN variant v ON v.id = pv.variant_id
LEFT JOIN product_size ps ON ps.product_id = p.id
LEFT JOIN size s ON s.id = ps.size_id
WHERE (p.name ILIKE '%' || $1 || '%' OR c.name ILIKE '%' || $1 || '%')
`

	args := []interface{}{search}
	argIndex := 2

	if len(categoryIDs) > 0 {
		query += fmt.Sprintf(" AND p.category_id = ANY($%d)", argIndex)
		args = append(args, categoryIDs)
		argIndex++
	}

	query += `
GROUP BY p.id, c.name
HAVING 1 = 1
`

	if minPrice != 0 {
		query += fmt.Sprintf(" AND COALESCE(MIN(ps.price), 0) >= $%d", argIndex)
		args = append(args, minPrice)
		argIndex++
	}
	if maxPrice != 0 {
		query += fmt.Sprintf(" AND COALESCE(MIN(ps.price), 0) <= $%d", argIndex)
		args = append(args, maxPrice)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY %s LIMIT %d OFFSET %d", orderByClause, limit, offset)

	rows, err := config.Db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var imagesJSON, variantsJSON, sizesJSON []byte

		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.MinPrice, &p.Stock,
			&p.Category, &imagesJSON, &variantsJSON, &sizesJSON,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		json.Unmarshal(imagesJSON, &p.Images)
		json.Unmarshal(variantsJSON, &p.Variants)
		json.Unmarshal(sizesJSON, &p.Sizes)

		products = append(products, p)
	}

	countQuery := `
	SELECT COUNT(*)
	FROM (
		SELECT p.id
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN product_size ps ON ps.product_id = p.id
		WHERE (p.name ILIKE '%' || $1 || '%' OR c.name ILIKE '%' || $1 || '%')
	`
	countArgs := []interface{}{search}
	countIndex := 2

	if len(categoryIDs) > 0 {
		countQuery += fmt.Sprintf(" AND p.category_id = ANY($%d)", countIndex)
		countArgs = append(countArgs, categoryIDs)
		countIndex++
	}

	countQuery += " GROUP BY p.id HAVING 1 = 1"

	if minPrice != 0 {
		countQuery += fmt.Sprintf(" AND COALESCE(MIN(ps.price), 0) >= $%d", countIndex)
		countArgs = append(countArgs, minPrice)
		countIndex++
	}

	if maxPrice != 0 {
		countQuery += fmt.Sprintf(" AND COALESCE(MIN(ps.price), 0) <= $%d", countIndex)
		countArgs = append(countArgs, maxPrice)
		countIndex++
	}

	countQuery += `) AS sub`

	var total int64
	err = config.Db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func GetProductByID(productID int64) (*Product, error) {
	ctx := context.Background()

	query := `
SELECT
  p.id,
  p.name,
  p.description,
  COALESCE(
    (SELECT MIN(price) FROM product_size WHERE product_id = p.id),
    p.price
  ) AS min_price,
  p.stock,
  COALESCE(c.name, '') AS category,

  COALESCE(
    json_agg(DISTINCT pi.image) FILTER (WHERE pi.image IS NOT NULL),
    '[]'
  ) AS images,

  COALESCE(json_agg(DISTINCT jsonb_build_object(
		'variant_id', v.id,
		'name', v.name
	  )
	) FILTER (WHERE v.id IS NOT NULL),
	'[]'
  ) AS variants,  

  COALESCE(
    json_agg(DISTINCT jsonb_build_object(
      'size_id', s.id,
      'size_name', s.name,
      'price', ps.price
    )) FILTER (WHERE s.id IS NOT NULL),
    '[]'
  ) AS sizes,

  p.created_at,
  p.updated_at

FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN product_img pi ON pi.product_id = p.id
LEFT JOIN product_variant pv ON pv.product_id = p.id
LEFT JOIN variant v ON v.id = pv.variant_id
LEFT JOIN product_size ps ON ps.product_id = p.id
LEFT JOIN size s ON s.id = ps.size_id

WHERE p.id = $1

GROUP BY p.id, c.name
LIMIT 1
`

	var p Product
	var imagesJSON, variantsJSON, sizesJSON []byte

	err := config.Db.QueryRow(ctx, query, productID).Scan(
		&p.ID, &p.Name, &p.Description, &p.MinPrice, &p.Stock,
		&p.Category, &imagesJSON, &variantsJSON, &sizesJSON,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(imagesJSON, &p.Images)
	_ = json.Unmarshal(variantsJSON, &p.Variants)
	_ = json.Unmarshal(sizesJSON, &p.Sizes)

	return &p, nil
}
func CreateProduct(req CreateProductRequest) (*Product, error) {
	ctx := context.Background()

	var productID int64
	err := config.Db.QueryRow(ctx,
		`INSERT INTO products (name, description, stock, category_id)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id`,
		req.Name, req.Description, req.Stock, req.CategoryID,
	).Scan(&productID)
	if err != nil {
		return nil, err
	}

	for _, img := range req.Images {
		_, _ = config.Db.Exec(ctx, `INSERT INTO product_img (image, product_id) VALUES ($1,$2)`,
			img, productID)
	}

	for _, v := range req.Variants {
		_, _ = config.Db.Exec(ctx, `INSERT INTO product_variant (variant_id, product_id) VALUES ($1,$2)`,
			v, productID)
	}

	for _, s := range req.Sizes {
		_, _ = config.Db.Exec(ctx, `INSERT INTO product_size (product_id, size_id, price) VALUES ($1,$2,$3)`,
			productID, s.SizeID, s.Price)
	}

	return GetProductByID(productID)
}

func UpdateProduct(id int64, req CreateProductRequest) (*Product, error) {
	ctx := context.Background()

	_, err := config.Db.Exec(ctx,
		`UPDATE products SET name=$1, description=$2, stock=$3, category_id=$4, updated_at=now()
		 WHERE id=$5`,
		req.Name, req.Description, req.Stock, req.CategoryID, id,
	)
	if err != nil {
		return nil, err
	}

	_, _ = config.Db.Exec(ctx, `DELETE FROM product_img WHERE product_id=$1`, id)
	for _, img := range req.Images {
		_, _ = config.Db.Exec(ctx, `INSERT INTO product_img (image, product_id) VALUES ($1,$2)`, img, id)
	}

	_, _ = config.Db.Exec(ctx, `DELETE FROM product_variant WHERE product_id=$1`, id)
	for _, v := range req.Variants {
		_, _ = config.Db.Exec(ctx, `INSERT INTO product_variant (variant_id, product_id) VALUES ($1,$2)`, v, id)
	}

	_, _ = config.Db.Exec(ctx, `DELETE FROM product_size WHERE product_id=$1`, id)
	for _, s := range req.Sizes {
		_, _ = config.Db.Exec(ctx,
			`INSERT INTO product_size (product_id, size_id, price) VALUES ($1,$2,$3)`,
			id, s.SizeID, s.Price,
		)
	}

	return GetProductByID(id)
}
func DeleteProduct(id int64) error {
    ctx := context.Background()

    _, err := config.Db.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
    if err != nil {
        return err
    }
    return nil
}

func UploadImgProduct(productID int64, imagePath string) error {
	ctx := context.Background()

	_, err := config.Db.Exec(ctx,
		`INSERT INTO product_img (image, product_id) VALUES ($1, $2)`,
		imagePath, productID,
	)
	return err
}

func GetRecommendationsByCategory(category string, excludeProductID int64) ([]Product, error) {
	ctx := context.Background()

	query := `
SELECT
  p.id,
  p.name,
  p.description,
  COALESCE(MIN(ps.price), 0) AS min_price,
  p.stock,
  COALESCE(c.name, '') AS category,

  COALESCE(
    json_agg(DISTINCT pi.image) FILTER (WHERE pi.image IS NOT NULL),
    '[]'
  ) AS images,

  COALESCE(
    json_agg(DISTINCT v.name) FILTER (WHERE v.name IS NOT NULL),
    '[]'
  ) AS variants,

  COALESCE(
    json_agg(DISTINCT jsonb_build_object(
      'size_id', s.id,
      'size_name', s.name,
      'price', ps.price
    )) FILTER (WHERE s.id IS NOT NULL),
    '[]'
  ) AS sizes,

  p.created_at,
  p.updated_at

FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN product_img pi ON pi.product_id = p.id
LEFT JOIN product_variant pv ON pv.product_id = p.id
LEFT JOIN variant v ON v.id = pv.variant_id
LEFT JOIN product_size ps ON ps.product_id = p.id
LEFT JOIN size s ON s.id = ps.size_id

WHERE c.name = $1 AND p.id != $2
GROUP BY p.id, c.name
ORDER BY p.created_at DESC
LIMIT 3
`

	rows, err := config.Db.Query(ctx, query, category, excludeProductID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		var imagesJSON, variantsJSON, sizesJSON []byte

		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.MinPrice, &p.Stock,
			&p.Category, &imagesJSON, &variantsJSON, &sizesJSON,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(imagesJSON, &p.Images)
		_ = json.Unmarshal(variantsJSON, &p.Variants)
		_ = json.Unmarshal(sizesJSON, &p.Sizes)

		products = append(products, p)
	}

	return products, nil
}
