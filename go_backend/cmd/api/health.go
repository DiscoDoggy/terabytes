package main

import (
	// "encoding/json"
	"net/http"
)

type Health struct {
	Status 	string 	`json:"status"`
	Env		string	`json:"environment"`
	Version	int		`json:"version"`
}
// HealthCheck godoc
//
//	@Summary		Fetches API health status
//	@Description	Fetches API health status	
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data:= map[string]string {
		"status": "ok",
		"env": app.config.env,
		"version": version,
	}
	
	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}