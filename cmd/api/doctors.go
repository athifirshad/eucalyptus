package main

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (app *application) getDoctorHandler(w http.ResponseWriter, r *http.Request) {
	doctorID, err := app.readIDParam(r)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	doctor, err := app.sqlc.GetDoctorById(r.Context(), int32(doctorID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "doctor not found"}, nil)
			return
		}
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"doctor": doctor}, nil)
}




func (app *application)getHealthRecordByRecordIdHandler(w http.ResponseWriter, r *http.Request) {
	// Read the recordID parameter using the utility function
	recordID, err := app.readIDParam(r)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	// Call the GetHealthRecordByRecordId function from your queries struct
	HealthRecord, err := app.sqlc.GetHealthRecordByRecordId(r.Context(), int32(recordID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "Health not found"}, nil)
			return
		}
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"HealthRecord": HealthRecord}, nil)
}

