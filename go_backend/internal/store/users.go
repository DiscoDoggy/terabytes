package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UsersStore struct {
	db *sql.DB
}

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password password `json:"-"`
	Created_at string `json:"created_at"`
	IsActivated bool `json:"is_activated"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	p.text = &password
	p.hash = bytes

	return nil
}

func (p *password) CheckPasswords(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func (s *UsersStore)CreateAndInvite(ctx context.Context, user *User, token string, expiry time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		err := s.Create(ctx, tx, user)
		if err != nil {
			return err
		}

		err = s.createUserInvite(ctx, tx, token, expiry, user.Id)
		if err != nil {
			return err
		}

		return nil
	})
	
	// For one transaction
		//create the user
		//create the email token 
}

func (s *UsersStore) createUserInvite(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userId string) error {
	insertInviteQuery := `
		INSERT INTO user_invitations (token, account_id, expiry)
		VALUES($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, insertInviteQuery, token, userId, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore)Create(ctx context.Context, tx *sql.Tx, user *User) error {
	insertUserQuery := `
		INSERT INTO accounts(username, email, password)
		VALUES($1, $2, $3)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, 
		insertUserQuery, 
		user.Username, 
		user.Email,
		user.Password.hash,
	).Scan(
		&user.Id,
		&user.Created_at,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint \"unique_accounts_username\""):
			return ErrUsernameConflict
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint \"unique_accounts_email\""):
			return ErrEmailConflict
		default:
			return err
		}
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

func (s *UsersStore) GetUserFeed(ctx context.Context, userId string, fq PaginatedFeedQuery) ([]FeedBlogPost, error) {
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
		LEFT JOIN blog_tags bt ON bt.blog_post_id = ub.id
		LEFT JOIN tags t ON t.id = bt.tag_id 
		WHERE f.user_id = $1
		GROUP BY ub.id, a.id
		ORDER BY ub.created_at DESC
		LIMIT $2
		OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, feedQuery, userId, fq.Limit, fq.Offset)
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
	defer rows.Close()

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

func(s *UsersStore) ActivateUser(ctx context.Context, hashToken string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, hashToken)
		if err != nil {
			return err
		}

		user.IsActivated = true
		err = s.updateUserActivation(ctx, tx, user)
		if err != nil {
			return err
		}

		err = s.cleanUserInvitations(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil

	})

}

func (s *UsersStore) updateUserActivation(ctx context.Context, tx *sql.Tx, user *User) error {
	updateUserActivationQuery := `
		UPDATE accounts 
		SET is_activated = true
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, updateUserActivationQuery, user.Id)
	if err != nil {
		return err
	}

	return nil
	
}

func (s *UsersStore) cleanUserInvitations(ctx context.Context, tx *sql.Tx, user *User) error {
	deleteInvitationQuery := `
		DELETE FROM user_invitations
		WHERE account_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, deleteInvitationQuery, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User,error) {
	getUserByTokenQuery := `
		SELECT 
			a.id, a.username, a.email, a.created_at, a.is_activated
		FROM accounts a
		JOIN user_invitations ui ON a.id = ui.account_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &User{}
	err := tx.QueryRowContext(
		ctx,
		getUserByTokenQuery, 
		token, 
		time.Now(),
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Created_at,
		&user.IsActivated,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
		
	}

	return user, nil

}