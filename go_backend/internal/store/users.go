package store

import (
	"context"
	"database/sql"

)

type UsersStore struct {
	db *sql.DB
}

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	Created_at string `json:"created_at"`
}

func (s *UsersStore)Create(ctx context.Context) error {
	return nil
}

func (s *UsersStore)GetUserById(ctx context.Context, userId string) (*User, error) {
	getUserQuery := `
		SELECT id, username, password, email, created_at
		FROM accounts
		WHERE id = $1
	`
	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User

	err := s.db.QueryRowContext(
		ctx, 
		getUserQuery, 
		userId,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Created_at,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
	
}

func (s *UsersStore) FollowUser(ctx context.Context, userId string, toFollowId string) error {
	//check if userId and ToFollowID exist as users
	insertFollowerQuery := `
		INSERT followers(user_id, follower_id)
		VALUES($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx, 
		insertFollowerQuery,
		userId,
		toFollowId,
	)
	if err != nil {
		return err
	}

	return nil
	
}

func (s *UsersStore) UnfollowUser(ctx context.Context, userId string, toUnfollowId string) error{
	doesUserFollowQuery := `
		SELECT user_id
		FROM followers
		WHERE user_id = $1 AND follower_id = $2
	`

	unFollowQuery := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var followId string
	err := s.db.QueryRowContext(
		ctx,
		doesUserFollowQuery,
		userId,
		toUnfollowId, 
	).Scan(
		&followId,
	)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	_, err = s.db.ExecContext(ctx, unFollowQuery, userId, toUnfollowId)
	if err != nil {
		return err
	}

	return nil
}