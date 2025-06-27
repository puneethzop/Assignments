package taskstore

import (
	"3layerarch/models"
	"database/sql"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) Create(task models.Task) error {
	_, err := s.DB.Exec("INSERT INTO TASKS (task, completed) VALUES (?, ?)", task.Task, task.Completed)
	return err
}

func (s *Store) GetByID(id int) (models.Task, error) {
	var t models.Task
	err := s.DB.QueryRow("SELECT id, task, completed FROM TASKS WHERE id = ?", id).
		Scan(&t.ID, &t.Task, &t.Completed)
	return t, err
}

func (s *Store) GetAll() ([]models.Task, error) {
	rows, err := s.DB.Query("SELECT id, task, completed FROM TASKS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Task, &t.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) Update(id int) error {
	_, err := s.DB.Exec("UPDATE TASKS SET completed = true WHERE id = ?", id)
	return err
}

func (s *Store) Delete(id int) error {
	_, err := s.DB.Exec("DELETE FROM TASKS WHERE id = ?", id)
	return err
}
