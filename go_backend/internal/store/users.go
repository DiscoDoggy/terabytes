package store

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (s *UsersStore)Create(ctx context.Context, user User) error {
	//check if user name or email is already registered
	userAlreadyExistsQuery := `
		SELECT COUNT(*)
		FROM accounts
		WHERE username = $1 OR email = $2
	`

	insertUserQuery := `
		INSERT INTO accounts(username, email, password)
		VALUES($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var numResults int
	err := s.db.QueryRowContext(ctx, userAlreadyExistsQuery, user.Username, user.Email).Scan(&numResults)
	if err != nil {
		return err
	}

	if numResults > 0 {
		return ErrConflict
	}

	_, err = s.db.ExecContext(ctx, 
		insertUserQuery, 
		user.Username, 
		user.Email,
		user.Password,
	)

	if err != nil {
		return err
	}

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

func (s *UsersStore) GetUserFeed(ctx context.Context, userId string) ([]FeedBlogPost, error) {
	feedQuery := `
		SELECT
			ub.id AS blog_post_id,
			a.id AS account_id,
			a.username,
			ub.title,
			ub.description,
			ub.created_at,
			COALESCE(
				JSONB_AGG(
					JSONB_BUILD_OBJECT('tag_id', bt.tag_id, 'tag_name', t."name")
				) FILTER (WHERE t.id IS NOT NULL), '[]'
			) AS tags
		FROM user_blogs ub
		JOIN accounts a ON a.id = ub.account_id
		JOIN followers f ON f.follower_id = a.id 
		LEFT JOIN tags t ON t.blog_id = ub.id 
		LEFT JOIN blog_tags bt ON bt.tag_id = t.id 
		WHERE f.user_id = $1
		ORDER BY ub.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, feedQuery, userId)
	if err != nil {
		return nil, err
	}
	feed := make([]FeedBlogPost, 0)
	for rows.Next() {
		var blogPost FeedBlogPost
		var feedTags string
		err := rows.Scan(
			&blogPost.Id,
			&blogPost.UserId,
			&blogPost.Username,
			&blogPost.Title,
			&blogPost.Description,
			&blogPost.CreatedAt,
			&feedTags,
		)
		if err != nil {
			return nil, err
		}

		var tags []Tag
		err = json.Unmarshal([]byte(feedTags), &tags)
		if err != nil {
			return nil, err
		}

		blogPost.Tags = tags

		feed = append(feed, blogPost)

	}

	return feed, nil
}