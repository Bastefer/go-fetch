package repository

import (
	"context"
	"fmt"
	"service-parser/internal/db/wrapper"
)

type BrandRepository interface {
	GetOrCreate(ctx context.Context,db wrapper.DB, name string) (int, error)
	GetCount(ctx context.Context,db wrapper.DB) (int, error)
}

type brandRepository struct {}

func NewBrandRepository() *brandRepository {
	return &brandRepository{}
}
func (r *brandRepository) GetOrCreate(
	ctx context.Context,
	db wrapper.DB,
	name string,
) (int, error) {
	const op = "internal/app/repository/brand/GetOrCreate"

	query := `
		INSERT INTO brands (name)
		VALUES ($1)
		ON CONFLICT (name)
		DO UPDATE SET
			name = EXCLUDED.name
		RETURNING id
	`

	var id int

	err := db.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *brandRepository) GetCount(ctx context.Context,db wrapper.DB) (int, error) {
	const op = "internal/app/repository/brand/GetCount"
	query := `
		SELECT count(*) from brands
	`
	var count int
	row := db.QueryRow(ctx, query)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
