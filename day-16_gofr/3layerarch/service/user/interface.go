package userservice

import "3layerarch/models"

type UserStore interface {
	CreateUser(u models.User) error
	GetUser(id int) (models.User, error)
}
