package store

import (
	"context"
	"database/sql"
)

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) FollowUser(ctx context.Context, userId string, toFollowId string) error {
	//check if userId and ToFollowID exist as users
	insertFollowerQuery := `
		INSERT INTO followers(user_id, follower_id)
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

func (s *FollowersStore) UnfollowUser(ctx context.Context, userId string, toUnfollowId string) error{
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