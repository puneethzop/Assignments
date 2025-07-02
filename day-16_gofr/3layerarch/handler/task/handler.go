package taskhandler

import (
	"strconv"

	"3layerarch/models"
	"gofr.dev/pkg/gofr"
)

type Handler struct {
	Service TaskService
}

func New(service TaskService) *Handler {
	return &Handler{Service: service}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a task with the given JSON body
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task to create"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Router /task [post]
func (h *Handler) CreateTask(ctx *gofr.Context) (interface{}, error) {
	var t models.Task

	err := ctx.Bind(&t)
	if err != nil {
		return nil, err
	}

	err = h.Service.CreateTask(ctx, t)

	if err != nil {
		return nil, err
	}

	return "Task created", nil
}

// GetTask godoc
// @Summary Get task by ID
// @Description Returns a task given its ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Bad Request"
// @Router /task/{id} [get]
func (h *Handler) GetTask(ctx *gofr.Context) (interface{}, error) {
	idStr := ctx.PathParam("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	t, err := h.Service.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// ViewTasks godoc
// @Summary List all tasks
// @Description Returns a list of all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Failure 500 {string} string "Internal Server Error"
// @Router /task [get]
func (h *Handler) ViewTasks(ctx *gofr.Context) (interface{}, error) {
	tasks, err := h.Service.ViewTasks(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// UpdateTask godoc
// @Summary Update task status
// @Description Marks a task as completed
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /task/{id} [put]
func (h *Handler) UpdateTask(ctx *gofr.Context) (interface{}, error) {
	idStr := ctx.PathParam("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	if err := h.Service.UpdateTask(ctx, id); err != nil {
		return nil, err
	}

	return "Task updated", nil
}

// DeleteTask godoc
// @Summary Delete task
// @Description Deletes a task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /task/{id} [delete]
func (h *Handler) DeleteTask(ctx *gofr.Context) (interface{}, error) {
	idStr := ctx.PathParam("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	if err := h.Service.DeleteTask(ctx, id); err != nil {
		return nil, err
	}

	return "Task deleted", nil
}
