package storage

import (
	"errors"

	"github.com/4aykovski/tg-notion-bot/internal/models"
)

var (
	ErrCantCreateNewPostgresDatabase = errors.New("can't create new postgres database")
	ErrCantCreateDatabaseConnection  = errors.New("can't create database connection")
	ErrCantPingDatabase              = errors.New("can't ping database")
	ErrCantCreateUser                = errors.New("can't create new user")
	ErrCantGetUser                   = errors.New("can't get user")
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id string) (*models.User, error)
}
