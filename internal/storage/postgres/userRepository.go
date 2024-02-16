package postgres

import (
	"context"
	"fmt"

	"github.com/4aykovski/tg-notion-bot/internal/models"
)

type UserRepository struct {
	db *Postgres
}

func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	q := `
		INSERT INTO users (id, name)
		VALUES ($1, $2)
	`

	_, err := repo.db.ExecContext(context.Background(), q, user.Id, user.Name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantCreateUser, err)
	}

	return user, nil
}

func (repo *UserRepository) GetUser(id string) (*models.User, error) {
	q := `
		SELECT id, name
		FROM users
		WHERE id = $1
	`

	rows, err := repo.db.QueryContext(context.Background(), q, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantGetUser, err)
	}

	var user *models.User
	if err = rows.Scan(user); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantGetUser, err)
	}

	return user, nil
}
