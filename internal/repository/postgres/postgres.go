package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	pool *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repo {
	return &repo{
		pool: db,
	}
}

func (r *repo) GetProducts(ctx context.Context, page int, limit int, sortOrder string) ([]models.Product, error) {
	//if limit > maxPaginationLimit {
	//	limit = maxPaginationLimit
	//}
	//
	//var sort string
	//switch sortOrder {
	//case "desc":
	//	sort = sortDescending
	//default:
	//	sort = sortAscending
	//}
	//
	//offset := (page - 1) * limit
	//
	//rows, err := r.pool.Query(ctx,
	//	`SELECT p.id, pp.name, pp.value
	//         FROM products p
	//         LEFT JOIN product_properties pp ON p.id = pp.product_id
	//         ORDER BY p.id `+sort+`
	//         LIMIT $1 OFFSET $2`, limit, offset)
	//
	//if err != nil {
	//	return nil, fmt.Errorf("productRepo - GetProducts - r.pool.Query: %w", err)
	//}
	//defer rows.Close()
	//
	//var products []models.Product
	//
	//for rows.Next() {
	//	var product models.Product
	//	var name, value *string
	//
	//	err := rows.Scan(&product.ID, &name, &value)
	//	if err != nil {
	//		return nil, fmt.Errorf("productRepo - GetProducts - rows.Scan: %w", err)
	//	}
	//	if product.Properties == nil {
	//		product.Properties = make(map[string]string)
	//	}
	//	if name != nil && value != nil {
	//		product.Properties[*name] = *value
	//	}
	//	products = append(products, product)
	//}
	//
	//if err := rows.Err(); err != nil {
	//	return nil, fmt.Errorf("productRepo - GetProducts - rows.Err: %w", err)
	//}
	//
	//return products, nil
}
