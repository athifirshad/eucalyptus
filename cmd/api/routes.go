package main

import (
	"net/http"
)

func (app *application) Routes() {
	app.router.Get("/", app.rootHandler)
	app.router.Get("/status", app.status)
	app.router.Get("/doctors/{id}", app.getDoctorHandler)
	app.router.Get("/HealthRecord/{id}", app.getHealthRecordByRecordIdHandler)
	app.router.Get("/GetHospitalByHospitalId/{id}", app.getHospitalByHospitalIdHandler)
	app.router.Post("/users", app.registerUserHandler)
	app.router.Put("/users/activated",app.activateUserHandler)
	app.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Page Not Found"))
	})
	app.router.Get("/docs", app.DocsHandler)

}
