package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
	All(ctx context.Context, userName string) ([]Page, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
}
