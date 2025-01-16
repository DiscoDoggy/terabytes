package main

import (
	"errors"
	"net/http"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
)

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload UserPayload
	err := readJSON(w, r, &userPayload)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	user := &store.User{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	}

	//hash password
	err = user.Password.HashPassword(userPayload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Users.Create(
		ctx,
		user,
	)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.resourceConflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
}