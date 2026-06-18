package connection


import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn *pgxpool.Pool
}

func NewStorage(ctx context.Context, dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{conn: conn}, nil
}

func (s *Storage) Pool() *pgxpool.Pool {
    return s.conn
}