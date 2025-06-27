package taskservice

import (
	"3layerarch/models"
	"3layerarch/store/task"
	"database/sql"
	"errors"
)

type Service struct {
	Store *taskstore.Store
}

func New(store *taskstore.Store) *Service {
	return &Service{Store: store}
}

func (s *Service) CreateTask(t models.Task) error {
	if t.Task == "" {
		return errors.New("task cannot be empty")
	}
	return s.Store.Create(t)
}

func (s *Service) GetTask(id int) (models.Task, error) {
	if id <= 0 {
		return models.Task{}, errors.New("invalid task ID")
	}
	task, err := s.Store.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, err
	}
	return task, nil
}

func (s *Service) ViewTasks() ([]models.Task, error) {
	return s.Store.GetAll()
}

func (s *Service) UpdateTask(id int) error {
	if id <= 0 {
		return errors.New("invalid task ID")
	}
	_, err := s.Store.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("task not found")
		}
		return err
	}
	return s.Store.Update(id)
}

func (s *Service) DeleteTask(id int) error {
	if id <= 0 {
		return errors.New("invalid task ID")
	}
	_, err := s.Store.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("task not found")
		}
		return err
	}
	return s.Store.Delete(id)
}
