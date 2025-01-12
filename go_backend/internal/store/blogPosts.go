package store

import (
	"context"
	"database/sql"
	// "errors"
	"strings"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/misc"
)

type BlogPostContent struct {
	Id 				string `json:"content_id"`
	BlogPostId 		string `json:"blog_post_id"` 
	ContentType 	string `json:"content_type"`
	ContentData 	string `json:"content_data"`
	ContentOrder 	int		`json:"content_order"`
}

type BlogPost struct {
	Id 		string `json:"id"`
	UserId	string `json:"user_id"`
	Username string `json:"username"`
	Title	string `json:"title"`
	Description string `json:"description"`
	Content []BlogPostContent `json:"content"`
	Tags 	[]string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

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

	createTagQuery := `
		INSERT INTO tags(name)
		VALUES($1)
		RETURNING id
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
		if err == sql.ErrNoRows {
			err = txn.QueryRowContext(
				ctx,
				createTagQuery,
				strings.TrimSpace(misc.CapitalizeString(blogPost.Tags[i])),
			).Scan(
				&currTagId,
			)

			if err != nil {
				_ = txn.Rollback()
				return err
			}
		}

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

func (s *BlogPostStore) GetBlogById(ctx context.Context, blogPostId string) (*BlogPost, error) {
	getBlogQuery := `
		SELECT 
			ub.id,
			a.username,
			a.id,
			ub.title,
			ub.description,
			ub.created_at,
			ubc.content_type,
			ubc.content_data,
			ubc.content_order,
			t."name",
			ubc.id,
			ub.updated_at
		FROM user_blogs ub
		JOIN accounts a ON a.id = ub.account_id
		JOIN user_blog_content ubc ON ub.id = ubc.user_blog_post_id
		LEFT JOIN blog_tags bt ON bt.blog_post_id = ub.id
		LEFT JOIN tags t ON t.id = bt.tag_id
		WHERE ub.id = $1 
		ORDER BY ubc.content_order ASC 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows,err := s.db.QueryContext(
		ctx, 
		getBlogQuery,
		blogPostId,
	)
	defer rows.Close()

	var blog BlogPost

	if err != nil {
		return nil, err
	}
	if !rows.Next()  {
		return nil, ErrNotFound
	}
	
	content := make([]BlogPostContent, 0)
	seenBlogContentIds := make(map[string]bool)
	tags := make([]string, 0)
	seenTags := make(map[string]bool)

	for rows.Next() {
		var tagNullStatus sql.NullString
		var currTag string
		var currContentId string
		var currContent BlogPostContent 
		err := rows.Scan(
			&blog.Id,
			&blog.Username,
			&blog.UserId,
			&blog.Title,
			&blog.Description,
			&blog.CreatedAt,
			&currContent.ContentType,
			&currContent.ContentData,
			&currContent.ContentOrder,
			&tagNullStatus,
			&currContentId,
			&blog.UpdatedAt,
		)
		if err != nil {
			return nil,err
		}

		if tagNullStatus.Valid {
			currTag = tagNullStatus.String
			
			if _, ok := seenTags[currTag]; !ok {
				seenTags[currTag] = true
				tags = append(tags, currTag)
			}

		}

		if _, ok := seenBlogContentIds[currContentId]; !ok {
			seenBlogContentIds[currContentId] = true
			currContent.Id = currContentId
			currContent.BlogPostId = blog.Id
			content = append(content, currContent)
		}

	}

	blog.Tags = tags
	blog.Content = content

	return &blog, nil
}

