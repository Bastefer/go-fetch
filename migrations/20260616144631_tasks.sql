-- +goose Up
SELECT 'up SQL query';
CREATE TABLE "tasks" (
    "id" SERIAL PRIMARY KEY,
    "started_at" timestamptz NOT NULL,
    "finished_at" timestamptz,
    "status" text NOT NULL
);
-- +goose Down
SELECT 'down SQL query';
DROP TABLE tasks;
