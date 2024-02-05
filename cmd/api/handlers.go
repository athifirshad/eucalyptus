package main

import (
	"fmt"
	"net/http"
)

func (app *application) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API!"))
}
func (app *application) hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", "1.0")
}