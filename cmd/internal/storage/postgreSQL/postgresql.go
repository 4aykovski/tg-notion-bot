package postgreSQL

import (
	"context"
	"database/sql"

	"github.com/4aykovski/tg-notion-bot/cmd/internal/storage"
	"github.com/4aykovski/tg-notion-bot/lib/helpers"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sql.Open("pq", dsn)
	if err != nil {
		return nil, helpers.ErrWrapIfNotNil("can't open database", err)
	}

	if err = db.Ping(); err != nil {
		return nil, helpers.ErrWrapIfNotNil("can't connect to database", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return helpers.ErrWrapIfNotNil("can't save a page", err)
	}

	return nil
}

func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? and user_name = ?`

	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return helpers.ErrWrapIfNotNil("can't delete a page", err)
	}

	return nil
}

func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? and user_name = ?`

	var count int

	err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count)
	if err != nil {
		return false, helpers.ErrWrapIfNotNil("can't check if page exists", err)
	}

	return count > 0, nil
}

func (s *Storage) All(ctx context.Context, userName string) ([]storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ?`

	var pages []storage.Page

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&pages)
	if err != nil {
		return nil, helpers.ErrWrapIfNotNil("can't get all pages", err)
	}

	if len(pages) < 0 {
		return nil, storage.ErrNoSavedPages
	}

	return pages, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return helpers.ErrWrapIfNotNil("can't init database", err)
	}

	return nil
}
