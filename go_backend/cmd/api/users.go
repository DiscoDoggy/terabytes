package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string
const userCtx userKey = "user"

type UserPayload struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}
// CreateUser godoc
//
//	@Summary		Creates a user
//	@Description	Creates a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UserPayload	true	"user payload"
//	@Success		200		{object}	string
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users [post]


// CreateUser godoc
//
//	@Summary		Fetches a user by id
//	@Description	Fetches a user by id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"user id"
//	@Success		200	{object}	string
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{user_id} [get]
func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {

	user := app.getUserFromCtx(r)

	err := writeJSON(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
// GetUserFeed godoc
//
//	@Summary		Fetches a user's feed
//	@Description	Fetches a user's feed
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery {
		Limit: 20,
		Offset: 0,
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(fq)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	
	user := app.getUserFromCtx(r)

	ctx := r.Context()

	userFeed, err := app.store.Users.GetUserFeed(ctx, user.Id, fq)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	err = writeJSON(w, http.StatusOK, userFeed)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
// FollowUser godoc
//
//	@Summary		Follows another user
//	@Description	Follows another user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"user id"
//	@Success		200	{object}	string
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{user_id}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	userToFollow := app.getUserFromCtx(r)
	tempCurrUser := "9efacfaf-2893-4665-b223-0ba333e04137" //TODO: CHANGE WHEN AUTH CAN FEED USER ID

	ctx := r.Context()

	err := app.store.Followers.FollowUser(ctx, tempCurrUser, userToFollow.Id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}
}

// FollowUser godoc
//
//	@Summary		Unfollows user
//	@Description	Unfollows user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"user id"
//	@Success		200	{object}	string
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{user_id}/unfollow [delete]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userToFollow := app.getUserFromCtx(r)
	tempCurrUser := "9efacfaf-2893-4665-b223-0ba333e04137" //TODO: CHANGE WHEN AUTH CAN FEED USER ID

	ctx := r.Context()

	err := app.store.Followers.UnfollowUser(ctx, tempCurrUser, userToFollow.Id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}
}

func (app *application) getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "user_id")

		ctx := r.Context()

		user, err := app.store.Users.GetUserById(ctx, userId)
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

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}