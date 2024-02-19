package postgres

import (
	"context"
	"fmt"

	"github.com/4aykovski/tg-notion-bot/internal/models"
	"github.com/4aykovski/tg-notion-bot/internal/storage"
)

type UserRepository struct {
	db *Postgres
}

func NewUserRepository(db *Postgres) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	q := `
		INSERT INTO "user" (id, name)
		VALUES ($1, $2)
	`

	_, err := repo.db.ExecContext(context.Background(), q, user.Id, user.Name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrCantCreateUser, err)
	}

	return user, nil
}

func (repo *UserRepository) GetUser(id int) (*models.User, error) {
	q := `
		SELECT id, name
		FROM "user"
		WHERE id = $1
	`

	row := repo.db.QueryRowContext(context.Background(), q, id)

	var user *models.User
	if err := row.Scan(user); err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrCantGetUser, err)
	}

	return user, nil
}
