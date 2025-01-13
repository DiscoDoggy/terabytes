package main

import (
	"net/http"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")
	
	ctx := r.Context()

	user, err := app.store.Users.GetUserById(ctx, userId)
	if err != nil {
		switch {
		case err == store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}