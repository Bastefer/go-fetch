package repository

import (
	"context"
	"fmt"
	"service-parser/internal/app/domain"
	"service-parser/internal/db/wrapper"
)

type ProductRepository interface {
	Upsert(ctx context.Context, db wrapper.DB,product domain.Product) error
	GetCount(ctx context.Context,db wrapper.DB) (int, error)
}

type productRepository struct {
}

func NewProductRepository() *productRepository {
	return &productRepository{
	}
}

func (r *productRepository) Upsert(
	ctx context.Context,
	db wrapper.DB,
	product domain.Product,
) error {
	const op = "internal/app/repository/product/Upsert"

	query := `
		INSERT INTO products (id, name, brand_id, category_id, price, stock)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name,
			brand_id = EXCLUDED.brand_id,
			category_id = EXCLUDED.category_id,
			price = EXCLUDED.price,
			stock = EXCLUDED.stock
	`

	_, err := db.Exec(
		ctx,
		query,
		product.ID,
		product.Name,
		product.BrandID,
		product.CategoryID,
		product.Price,
		product.Stock,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *productRepository) GetCount(ctx context.Context,db wrapper.DB) (int, error) {
	const op = "internal/app/repository/product/GetCount"
	query := `
		SELECT count(*) from products
	`
	var count int
	row := db.QueryRow(ctx, query)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
