package main

import (
	"net/http"
	"errors"
	"context"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string
const userCtx userKey = "user"

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {

	user := app.getUserFromCtx(r)

	err := writeJSON(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

}

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