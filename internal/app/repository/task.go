package repository

import (
	"context"
	"errors"
	"fmt"
	"service-parser/internal/app/domain"
	"service-parser/internal/db/wrapper"
	"time"

	"github.com/jackc/pgx/v5"
)

type TaskRepository interface {
	Create(ctx context.Context,db wrapper.DB) (int, error)

	GetRunning(ctx context.Context, db wrapper.DB) (*domain.DownloadTask, error)

	UpdateStatus(
		ctx context.Context,
		db wrapper.DB,
		id int,
		status domain.TaskStatus,
		finishedAt *time.Time,
	) error
}

type taskRepository struct {
}

func NewTaskRepository() *taskRepository {
	return &taskRepository{
	}
}
func (r *taskRepository) Create(ctx context.Context,db wrapper.DB) (int, error) {
	const op = "internal/app/repository/task/Create"
	var id int
	query := `
		INSERT INTO tasks (started_at, status)
		VALUES (NOW(), $1)
		RETURNING id
	`
	err := db.QueryRow(ctx, query, domain.TaskRunning).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil

}
func (r *taskRepository) GetRunning(ctx context.Context,db wrapper.DB) (*domain.DownloadTask, error) {
	const op = "internal/app/repository/task/GetRunning"

	var t domain.DownloadTask

	query := `SELECT id, started_at, finished_at, status
		FROM tasks
		WHERE status = $1
		ORDER BY started_at DESC
		LIMIT 1`

	err := db.QueryRow(ctx, query, domain.TaskRunning).Scan(
		&t.ID,
		&t.StartedAt,
		&t.FinishedAt,
		&t.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &t, nil
}
func (r *taskRepository) UpdateStatus(
	ctx context.Context,
	db wrapper.DB,
	id int,
	status domain.TaskStatus,
	finishedAt *time.Time,
) error {
	const op = "internal/app/repository/task/UpdateStatus"
	_, err := db.Exec(ctx, `
		UPDATE tasks
		SET status = $1,
		    finished_at = $2
		WHERE id = $3
	`, status, finishedAt, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
