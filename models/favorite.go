package models

import (
	"backend/config"
	"context"
	"encoding/json"
)

type FavoriteReq struct {
	Id          int    `json:"id"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

func Favorite(page, limit int) ([]Product, int64, error) {
	ctx := context.Background()

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 4
	}

	offset := (page - 1) * limit

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

WHERE p.is_favorite = TRUE

GROUP BY p.id, c.name
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2
`

	rows, err := config.Db.Query(ctx, query, limit, offset)
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
		WHERE p.is_favorite = TRUE`,
	).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
