package taskservice

import (
	//"3layerarch/handler/userhandler"
	"3layerarch/models"
	"database/sql"
	"errors"
	"fmt"
)

type TaskStore interface {
	CreateTask(t models.Task) error
	GetTask(id int) (models.Task, error)
	ViewTasks() ([]models.Task, error)
	UpdateTask(id int) error
	DeleteTask(id int) error
}

type UserService interface {
	GetUser(id int) (models.User, error)
}

type Service struct {
	TaskStore   TaskStore
	UserService UserService
}

func New(ts TaskStore, us UserService) *Service {
	return &Service{TaskStore: ts, UserService: us}
}

func (s *Service) CreateTask(t models.Task) error {
	fmt.Println("CreateTask received:", t)

	if t.Task == "" {
		return errors.New("task cannot be empty")
	}
	// Validate user existence before creating task
	_, err := s.UserService.GetUser(t.UserID)
	if err != nil {
		return errors.New("user ID not found")
	}
	return s.TaskStore.CreateTask(t)
}

func (s *Service) GetTask(id int) (models.Task, error) {
	if id <= 0 {
		return models.Task{}, errors.New("invalid task ID")
	}
	task, err := s.TaskStore.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, err
	}
	return task, nil
}

func (s *Service) ViewTasks() ([]models.Task, error) {
	return s.TaskStore.ViewTasks()
}

func (s *Service) UpdateTask(id int) error {
	if id <= 0 {
		return errors.New("invalid task ID")
	}
	_, err := s.TaskStore.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("task not found")
		}
		return err
	}
	return s.TaskStore.UpdateTask(id)
}

func (s *Service) DeleteTask(id int) error {
	if id <= 0 {
		return errors.New("invalid task ID")
	}
	_, err := s.TaskStore.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("task not found")
		}
		return err
	}
	return s.TaskStore.DeleteTask(id)
}
