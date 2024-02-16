package storage

import "github.com/4aykovski/tg-notion-bot/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id string) (*models.User, error)
}
