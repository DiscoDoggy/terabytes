package store

import (
	"context"
	"database/sql"
	"time"
)

var QueryTimeoutDuration = time.Second * 5

type Storage struct {
	Posts interface {
		Create(context.Context, *BlogPost) error
		GetBlogById(context.Context, string) (BlogPost, error)
	}

	Users interface {
		Create(context.Context) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage {
		Posts: &BlogPostStore{db},
		Users: &UsersStore {db},

	}
}
