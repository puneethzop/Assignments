package taskstore

import (
	"3layerarch/models"
	"database/sql"
	"log"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTask(t models.Task) error {
	_, err := s.db.Exec("INSERT INTO TASKS (task, completed, user_id) VALUES (?, ?, ?)", t.Task, t.Completed, t.UserID)
	return err
}

func (s *Store) GetTask(id int) (models.Task, error) {
	var t models.Task
	err := s.db.QueryRow("SELECT id, task, completed, user_id FROM TASKS WHERE id = ?", id).
		Scan(&t.ID, &t.Task, &t.Completed, &t.UserID)
	return t, err
}

func (s *Store) ViewTasks() ([]models.Task, error) {
	rows, err := s.db.Query("SELECT id, task, completed, user_id FROM TASKS")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Error closing rows:", err)
		}
	}()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Task, &t.Completed, &t.UserID); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) UpdateTask(id int) error {
	_, err := s.db.Exec("UPDATE TASKS SET completed = true WHERE id = ?", id)
	return err
}

func (s *Store) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM TASKS WHERE id = ?", id)
	return err
}
