package userhandler

import "3layerarch/models"

type UserService interface {
	CreateUser(u models.User) error
	GetUser(id int) (models.User, error)
}
