package store

import (
	"context"
	"database/sql"
	"strings"
)

type BlogPostContent struct {
	Id 				string
	BlogPostId 		string
	ContentType 	string
	ContentData 	string
	ContentOrder 	int
}

type BlogPost struct {
	Id 		string
	UserId	string
	Title	string
	Description string
	Content []BlogPostContent
	Tags 	[]string
	CreatedAt string
	UpdatedAt string

}

type BlogPostStore struct {
	db *sql.DB
}

func (s *BlogPostStore) Create(ctx context.Context, blogPost * BlogPost) error {
	insertPostInfoQuery := `
		INSERT INTO user_blogs(account_id, title, description)
		VALUES($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	insertPostContentQuery := `
		INSERT INTO user_blog_content(user_blog_post_id, content_type, content_data, content_order)
		VALUES($1, $2, $3, $4)
	`

	getTagIdQuery := `
		SELECT id FROM tags WHERE lower(name) = lower($1)
	`

	insertPostTagsQuery := `
		INSERT INTO blog_tags(blog_post_id, tag_id)
		VALUES($1, $2)
	`
	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	txn, err := s.db.Begin()
	if err != nil {
		_ = txn.Rollback()
		return err
	}
	
	err = txn.QueryRowContext(
		ctx,
		insertPostInfoQuery,
		blogPost.UserId,
		blogPost.Title,
		blogPost.Description,
	).Scan(
		&blogPost.Id,
		&blogPost.CreatedAt,
		&blogPost.UpdatedAt,
	)
	if err != nil {
		_ = txn.Rollback()
		return err
	}

	var insertedBlogId string = blogPost.Id
	for i := 0; i < len(blogPost.Content); i++ {
		_, err := txn.Exec(
			insertPostContentQuery, 
			insertedBlogId, 
			blogPost.Content[i].ContentType,
			blogPost.Content[i].ContentData,
			blogPost.Content[i].ContentOrder,
		)
		if err != nil{
			_ = txn.Rollback()
			return err
		}
	}

	for i := 0; i < len(blogPost.Tags); i++ {
		var currTagId string
		err := txn.QueryRowContext(
			ctx,
			getTagIdQuery,
			strings.TrimSpace(blogPost.Tags[i]),
		).Scan(
			&currTagId,
		)
		if err != nil {
			_ = txn.Rollback()
			return err
		}

		_, err = txn.Exec(insertPostTagsQuery, insertedBlogId, currTagId)
		if err != nil {
			_ = txn.Rollback()
			return err
		}
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
