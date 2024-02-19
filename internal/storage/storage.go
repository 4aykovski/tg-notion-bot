package storage

import (
	"errors"

	"github.com/4aykovski/tg-notion-bot/internal/models"
)

var (
	ErrCantCreateNewDatabase        = errors.New("can't create new database")
	ErrCantCreateDatabaseConnection = errors.New("can't create database connection")
	ErrCantPingDatabase             = errors.New("can't ping database")
	ErrCantCreateUser               = errors.New("can't create new user")
	ErrCantGetUser                  = errors.New("can't get user")
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id int) (*models.User, error)
}
