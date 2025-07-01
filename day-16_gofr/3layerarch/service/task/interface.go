package taskservice

import "3layerarch/models"

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
