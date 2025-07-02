package taskservice

import (
	"3layerarch/models"
	"gofr.dev/pkg/gofr"
)

type TaskStore interface {
	CreateTask(ctx *gofr.Context, t models.Task) error
	GetTask(ctx *gofr.Context, id int) (models.Task, error)
	ViewTasks(ctx *gofr.Context) ([]models.Task, error)
	UpdateTask(ctx *gofr.Context, id int) error
	DeleteTask(ctx *gofr.Context, id int) error
}

type UserService interface {
	GetUser(ctx *gofr.Context, id int) (models.User, error)
}
