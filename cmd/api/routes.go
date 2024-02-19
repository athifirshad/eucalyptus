package main

import (
	"net/http"
	"os"

	"github.com/go-chi/docgen"
)

func (app *application) Routes() {
	app.router.Get("/", app.rootHandler)
	app.router.Get("/status", app.status)
	app.router.Get("/doctors/{id}", app.getDoctorHandler)
	app.router.Post("/users", app.registerUserHandler)
	app.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Page Not Found"))
	})
	app.router.Get("/docs", app.DocsHandler)

}

func (app *application) DocsHandler(w http.ResponseWriter, r *http.Request) {
	docs := docgen.MarkdownRoutesDoc(app.router, docgen.MarkdownOpts{
		ProjectPath: "github.com/athifirshad/eucalyptus",
		// Intro text included at the top of the generated markdown file.
		Intro: "Generated documentation for Eucalyptus",
	})
	err := os.WriteFile("./docs/routes.md", []byte(docs), 0644)
	if err != nil {
		app.logger.Sugar().Errorf("error writing docs", "err", err)
		http.Error(w, "Internal server error fetching docs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Documentation successfully written to routes.md"))
}
