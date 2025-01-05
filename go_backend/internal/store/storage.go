package store

import (
	"database/sql"
	"context"
)

type Storage struct {
	Posts interface {
		Create(context.Context) error
	}

	Users interface {
		Create(context.Context) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage {
		Posts: &PostsStore{db},
		Users: &UsersStore {db},

	}
}