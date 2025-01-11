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

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data:= map[string]string {
		"status": "ok",
		"env": app.config.env,
		"version": version,
	}
	
	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}