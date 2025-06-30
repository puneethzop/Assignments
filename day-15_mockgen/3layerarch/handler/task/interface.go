package taskhandler

import "3layerarch/models"

type TaskService interface {
	CreateTask(t models.Task) error
	GetTask(id int) (models.Task, error)
	ViewTasks() ([]models.Task, error)
	UpdateTask(id int) error
	DeleteTask(id int) error
}
