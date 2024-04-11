package main

import (
	"net/http"

	"github.com/go-chi/cors"
)

func (app *application) Routes() {
	app.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, 
	}))

	app.router.Get("/", app.rootHandler)
	app.router.Get("/status", app.status)
	app.router.Get("/doctors/{id}", app.requireActivatedUser(app.getDoctorHandler))
	app.router.Get("/HealthRecord/{id}", app.requireActivatedUser(app.getHealthRecordByRecordIdHandler))
	app.router.Get("/GetHospitalByHospitalId/{id}", app.requireActivatedUser(app.getHospitalByHospitalIdHandler))
	app.router.Get("/GetTreatmentHistoryByPatientID/{id}", app.requireActivatedUser(app.GetTreatmentHistoryByPatientIDHandler))

	app.router.Post("/users", app.registerUserHandler)
	app.router.Put("/users/activated", app.activateUserHandler)
	app.router.Post("/tokens/authentication",app.createAuthenticationTokenHandler)
	app.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Page Not Found"))
	})
	app.router.Get("/docs", app.DocsHandler)

}
