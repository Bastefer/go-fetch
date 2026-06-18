package repository

import (
	"context"
	"fmt"
	"service-parser/internal/app/domain"
	"service-parser/internal/db/wrapper"

)

type ClientRepository interface {
	Upsert(ctx context.Context, db wrapper.DB,client domain.Client) error

	GetCount(ctx context.Context,db wrapper.DB) (int, error)
}

type clientRepository struct {
}

func NewClientRepository() *clientRepository {
	return &clientRepository{
	}
}

func (r *clientRepository) Upsert(ctx context.Context,db wrapper.DB, client domain.Client) error {
	const op = "internal/app/repository/client/Upsert"

	query := `
		INSERT INTO clients (id, first_name, last_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (id)
		DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name
	`

	_, err := db.Exec(
		ctx,
		query,
		client.ID,
		client.FirstName,
		client.LastName,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *clientRepository) GetCount(ctx context.Context,db wrapper.DB) (int, error) {
	const op = "internal/app/repository/client/GetCount"
	query := `
		SELECT count(*) from clients
	`
	var count int
	row := db.QueryRow(ctx, query)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
