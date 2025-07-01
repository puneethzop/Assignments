package userhandler

import (
	"strconv"

	"3layerarch/models"
	"gofr.dev/pkg/gofr"
)

type Handler struct {
	Service UserService
}

func New(service UserService) *Handler {
	return &Handler{Service: service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a user with the given JSON body
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User to create"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Router /user [post]
func (h *Handler) CreateUser(ctx *gofr.Context) (interface{}, error) {
	var u models.User
	if err := ctx.Bind(&u); err != nil {
		return nil, err
	}

	if err := h.Service.CreateUser(u); err != nil {
		return nil, err
	}
	
	return "error", nil
}

// GetUser godoc
// @Summary Get user by ID
// @Description Returns a user given their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Bad Request"
// @Router /user/{id} [get]
func (h *Handler) GetUser(ctx *gofr.Context) (interface{}, error) {
	idStr := ctx.PathParam("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	u, err := h.Service.GetUser(id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
