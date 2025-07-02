package taskstore

import (
	"3layerarch/models"
	"database/sql"
	"gofr.dev/pkg/gofr"
	"log"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTask(ctx *gofr.Context, t models.Task) error {
	_, err := ctx.SQL.ExecContext(ctx, "INSERT INTO TASKS (task, completed, user_id) VALUES (?, ?, ?)", t.Task, t.Completed, t.UserID)
	return err
}

func (s *Store) GetTask(ctx *gofr.Context, id int) (models.Task, error) {
	var t models.Task
	err := ctx.SQL.QueryRowContext(ctx, "SELECT id, task, completed, user_id FROM TASKS WHERE id = ?", id).
		Scan(&t.ID, &t.Task, &t.Completed, &t.UserID)
	return t, err
}

// In task.go
func (s *Store) ViewTasks(ctx *gofr.Context) ([]models.Task, error) {
	rows, err := ctx.SQL.QueryContext(ctx, "SELECT id, task, completed, user_id FROM TASKS")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Error closing rows:", err)
		}
	}()

	tasks := make([]models.Task, 0) // CHANGE THIS: Initialize as an empty, non-nil slice
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Task, &t.Completed, &t.UserID); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
func (s *Store) UpdateTask(ctx *gofr.Context, id int) error {
	_, err := ctx.SQL.ExecContext(ctx, "UPDATE TASKS SET completed = true WHERE id = ?", id)
	return err
}

func (s *Store) DeleteTask(ctx *gofr.Context, id int) error {
	_, err := ctx.SQL.ExecContext(ctx, "DELETE FROM TASKS WHERE id = ?", id)
	return err
}
