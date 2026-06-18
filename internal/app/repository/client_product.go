package repository

import (
	"context"
	"fmt"
	"service-parser/internal/db/wrapper"

)

type ClientProductRepository interface {
	ReplaceProducts(
		ctx context.Context,
		db wrapper.DB,
		clientID int,
		productIDs []int,
	) error
}

type clientProductRepository struct {
}

func NewClientProductRepository() *clientProductRepository {
	return &clientProductRepository{
	}
}
func (r *clientProductRepository) ReplaceProducts(
	ctx context.Context,
	db wrapper.DB,
	clientID int,
	productIDs []int,
) error {
	const op = "internal/app/repository/client_product/ReplaceProducts"

	query := `DELETE FROM clients_products WHERE client_id = $1`

	_, err := db.Exec(
		ctx,
		query,
		clientID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, productID := range productIDs {
		insertQuery := `
				INSERT INTO clients_products (client_id, product_id)
				VALUES ($1, $2)
			`
		_, err := db.Exec(
			ctx,
			insertQuery,
			clientID,
			productID,
		)
		if err != nil {
			return fmt.Errorf("%s: insert: %w", op, err)
		}
	}

	return nil
}
