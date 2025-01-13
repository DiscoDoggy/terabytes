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