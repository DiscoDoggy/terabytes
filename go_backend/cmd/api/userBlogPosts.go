package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type blogKey string
const blogCtx blogKey = "blog"

type BlogPostContentPayload struct {
	ContentType 	string 	`json:"content_type" validate:"required"`
	ContentData 	string	`json:"content_data" validate:"required"`
	ContentOrder 	int		`json:"content_order" validate:"required"`
}

type CreateBlogPostPayload struct {
	Title 				string `json:"title" validate:"required,max=300"`
	Description 		string `json:"description" validate:"required,max=2200"`
	Content 			[]BlogPostContentPayload `json:"content" validate:"required"`
	Tags 				[]string `json:"tags"`
}

// CreateBlogPost godoc
//
//	@Summary		Creates a user blog post
//	@Description	Creates a user blog post
//	@Tags			user_blogs
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateBlogPostPayload	true	"BlogPost payload"
//	@Success		200		{object}	string
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/blogs [post]
func (app *application) createBlogHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateBlogPostPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	blogPostContent := make([]store.BlogPostContent, len(payload.Content))
	for i:=0; i < len(blogPostContent); i++ {
		blogPostContent[i].ContentType = payload.Content[i].ContentType
		blogPostContent[i].ContentData = payload.Content[i].ContentData
		blogPostContent[i].ContentOrder = payload.Content[i].ContentOrder
	}

	blogPost := &store.BlogPost{
		UserId: "9efacfaf-2893-4665-b223-0ba333e04137",
		Username: "",
		Title: payload.Title,
		Description: payload.Description,
		Content: blogPostContent,
		Tags: payload.Tags,
	}

	ctx := r.Context()
	
	err = app.store.Posts.Create(ctx, blogPost)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// GetBlogById godoc
//
//	@Summary		Fetches a user blog post
//	@Description	Fetches a user blog post
//	@Tags			user_blogs
//	@Accept			json
//	@Produce		json
//	@Param			blog_id	path		string	true	"blog post ID"
//	@Success		200		{object}	store.BlogPost
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/blogs/{blog_id} [get]
func (app *application) getBlogByIdHandler(w http.ResponseWriter, r *http.Request) {
	blog := app.getBlogFromCtx(r)

	fmt.Println(len(blog.Content))

	err := writeJSON(w, http.StatusOK, blog)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}


// DeleteBlogById godoc
//
//	@Summary		Deletes a user blog post
//	@Description	Deletes a user blog post
//	@Tags			user_blogs
//	@Accept			json
//	@Produce		json
//	@Param			blog_id	path		string	true	"blog post ID"
//	@Success		200		{object}	string
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/blogs/{blog_id} [delete]
func (app *application) deleteBlogByIdHandler(w http.ResponseWriter, r *http.Request) {
	blogId := chi.URLParam(r, "blog_id")

	ctx := r.Context()

	err := app.store.Posts.DeleteBlogById(ctx, blogId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}
}

func (app *application) getBlogFromCtx(r *http.Request) *store.BlogPost {
	blog, _ := r.Context().Value(blogCtx).(*store.BlogPost)
	return blog
}

func (app *application) blogPostContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		blogId := chi.URLParam(r, "blog_id")

		ctx := r.Context()
		blog, err := app.store.Posts.GetBlogById(ctx, blogId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r , err)
				return
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, blogCtx, blog)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

