package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
	"github.com/google/uuid"
)

// CreateUser godoc
//
//	@Summary		Creates a user
//	@Description	Creates a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UserPayload	true	"user credentials"
//	@Success		200		{object}	string
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/authentication/user [post]
func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	//create user
	var userPayload UserPayload
	err := readJSON(w, r, &userPayload)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()
	//store token in DB hash
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	fmt.Println("PLAIN TOKEN:", plainToken)
	app.logger.Infow("PLAIN TOKEN", "token", plainToken)

	user := &store.User{
		Username: userPayload.Username,
		PhoneNumber: userPayload.PhoneNumber,
		Email:    userPayload.Email,
	}

	//hash password
	err = user.Password.HashPassword(userPayload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Users.CreateAndInvite(
		ctx,
		user,
		hashToken,
		app.config.mail.exp,
	)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrUsernameConflict):
			app.badRequestError(w, r, err)
			return
		case errors.Is(err, store.ErrEmailConflict):
			app.badRequestError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	//send email
}