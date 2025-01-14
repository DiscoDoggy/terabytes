package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)
var ErrNotFound = errors.New("record not found")
var QueryTimeoutDuration = time.Second * 5

type Storage struct {
	Posts interface {
		Create(context.Context, *BlogPost) error
		GetBlogById(context.Context, string) (*BlogPost, error)
		DeleteBlogById(context.Context, string) error
	}

	Users interface {
		Create(context.Context) error
		GetUserById(context.Context, string) (*User, error)
		FollowUser(context.Context, string, string) error
		UnfollowUser(context.Context, string, string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage {
		Posts: &BlogPostStore{db},
		Users: &UsersStore {db},

	}
}
