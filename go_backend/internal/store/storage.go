package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)
var ErrNotFound = errors.New("record not found")
var ErrConflict = errors.New("input information conflicts with already existing information")
var ErrUsernameConflict = errors.New("inputted username is already taken")
var ErrEmailConflict = errors.New("inputted email is already taken")
var QueryTimeoutDuration = time.Second * 5

type Storage struct {
	Posts interface {
		Create(context.Context, *BlogPost) error
		GetBlogById(context.Context, string) (*BlogPost, error)
		DeleteBlogById(context.Context, string) error
	}

	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		createUserInvite(context.Context, *sql.Tx, string, time.Duration, string) error
		GetUserById(context.Context, string) (*User, error)
		GetUserFeed(context.Context, string, PaginatedFeedQuery) ([]FeedBlogPost, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		ActivateUser(context.Context, string) error
		
	}

	Followers interface {
		FollowUser(context.Context, string, string) error
		UnfollowUser(context.Context, string, string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage {
		Posts: &BlogPostStore{db},
		Users: &UsersStore{db},
		Followers: &FollowersStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error ) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
