package postgres

import (
	"database/sql"
	"fmt"

	"github.com/4aykovski/tg-notion-bot/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	*sql.DB
}

func NewPostgresDatabase(cfg config.DatabaseConfig) (*Postgres, error) {
	db, err := NewConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantCreateNewPostgresDatabase, err)
	}

	return &Postgres{db}, nil
}

func NewConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.DSNTemplate

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantCreateDatabaseConnection, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantPingDatabase, err)
	}

	return db, nil
}
