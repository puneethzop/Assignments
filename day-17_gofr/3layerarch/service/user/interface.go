package userservice

import (
	"3layerarch/models"
	"gofr.dev/pkg/gofr"
)

type UserStore interface {
	CreateUser(ctx *gofr.Context, u models.User) error
	GetUser(ctx *gofr.Context, id int) (models.User, error)
}
