package models

import (
	"backend/config"
	"context"
	"encoding/json"
	"time"
)

type ProductSize struct {
	SizeID   int64   `json:"size_id" binding:"required,gt=0"`
	SizeName string  `json:"size_name"`
	Price    float64 `json:"price" binding:"required,gte=0"`
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

func GetProducts(page, limit int, search, sort string) ([]Product, int64, error) {
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
	COALESCE(MIN(ps.price), 0) AS min_price,
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

WHERE 
	p.name ILIKE '%' || $1 || '%'
	OR c.name ILIKE '%' || $1 || '%'

GROUP BY p.id, c.name

ORDER BY ` + orderByClause + `
LIMIT $2 OFFSET $3
`

	rows, err := config.Db.Query(ctx, query, search, limit, offset)
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

	var total int64
	err = config.Db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE (p.name ILIKE '%' || $1 || '%' OR c.name ILIKE '%' || $1 || '%')`,
		search,
	).Scan(&total)

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

	_, err := config.Db.Exec(ctx, `DELETE FROM product_img WHERE product_id=$1`, id)
	if err != nil {
		return err
	}

	_, err = config.Db.Exec(ctx, `DELETE FROM product_variant WHERE product_id=$1`, id)
	if err != nil {
		return err
	}

	_, err = config.Db.Exec(ctx, `DELETE FROM product_size WHERE product_id=$1`, id)
	if err != nil {
		return err
	}

	_, err = config.Db.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
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
