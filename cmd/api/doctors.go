package main

import "net/http"

func (app *application) getDoctorHandler(w http.ResponseWriter, r *http.Request) {
	doctorID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	doctor, err := app.sqlc.GetDoctorById(r.Context(), int32(doctorID))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"doctor": doctor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
