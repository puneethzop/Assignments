package taskhandler

import (
	"3layerarch/models"
	"gofr.dev/pkg/gofr"
)

type TaskService interface {
	CreateTask(ctx *gofr.Context, t models.Task) error
	GetTask(ctx *gofr.Context, id int) (models.Task, error)
	ViewTasks(ctx *gofr.Context) ([]models.Task, error)
	UpdateTask(ctx *gofr.Context, id int) error
	DeleteTask(ctx *gofr.Context, id int) error
}
