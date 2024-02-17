package main

import (
	"net/http"
)

func (app *application) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API!"))
}
func (app *application) status(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     "1.0",
		},
	}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
