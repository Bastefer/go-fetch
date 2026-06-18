package repository

import (
	"context"
	"fmt"
	"service-parser/internal/db/wrapper"
)

type CategoryRepository interface {
	GetOrCreate(ctx context.Context,db wrapper.DB, name string) (int, error)
	GetCount(ctx context.Context,db wrapper.DB) (int, error)
}

type categoryRepository struct {}

func NewCategoryRepository() *categoryRepository {
	return &categoryRepository{
	}
}
func (r *categoryRepository) GetOrCreate(ctx context.Context,db wrapper.DB, name string) (int, error) {
	const op = "internal/app/repository/category/GetOrCreate"

	query := `
		INSERT INTO categories (name)
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

func (r *categoryRepository) GetCount(ctx context.Context,db wrapper.DB) (int, error) {
	const op = "internal/app/repository/category/GetCount"
	query := `
		SELECT count(*) from categories
	`
	var count int
	row := db.QueryRow(ctx, query)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
