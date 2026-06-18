package domain

import "time"

type TaskStatus string

const (
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
)

type DownloadTask struct {
	ID         int
	StartedAt  time.Time
	FinishedAt *time.Time
	Status     TaskStatus
}
