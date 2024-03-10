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

	doctor, err := app.sqlc.GetDoctorById(r.Context(), doctorID)
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

func (app *application) getHealthRecordByRecordIdHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) getHealthRecordsByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	// Read the patientID parameter using the utility function
	patientID, err := app.readIDParam(r)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}
	// Call the GetHealthRecordsByPatientId function from your queries struct
	healthRecords, err := app.sqlc.GetHealthRecordsByPatientId(r.Context(), patientID)
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"healthRecords": healthRecords}, nil)
}

func (app *application) getHospitalByHospitalIdHandler(w http.ResponseWriter, r *http.Request) {
	// Read the hospitalID parameter using the utility function
	hospitalID, err := app.readIDParam(r)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	// Call the GetHospitalByHospitalId function from your queries struct
	hospital, err := app.sqlc.GetHospitalByHospitalId(r.Context(), hospitalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "Hospital not found"}, nil)
			return
		}
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"hospital": hospital}, nil)
}

func (app *application) getMedicationsByPrescriptionIdHandler(w http.ResponseWriter, r *http.Request) {
	prescriptionID, err := app.readIDParam(r)
	// Read the prescriptionID parameter using the utility function
	medications, err := app.sqlc.GetMedicationsByPrescriptionId(r.Context(), prescriptionID)
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"medications": medications}, nil)
}

func (app *application) getPrescriptionsByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	// Read the patientID parameter using the utility function
	patientID, err := app.readIDParam(r)

	prescriptions, err := app.sqlc.GetPrescriptionsByPatientId(r.Context(), patientID)
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"prescriptions": prescriptions}, nil)
}

func (app *application) getProfileByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	// Read the userID parameter using the utility function
	userID, err := app.readIDParam(r)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	// Call the GetProfileByUserId function from your queries struct
	profile, err := app.sqlc.GetProfileByUserId(r.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "Profile not found"}, nil)
			return
		}
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
}

func (app *application) GetTreatmentHistoryByPatientIDHandler(w http.ResponseWriter, r *http.Request) {
	// Read the patientID parameter using the utility function
	patientID, err := app.readIDParam(r)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	// Call the GetTreatmentHistoryByPatientID function from your queries struct
	treatmentHistory, err := app.sqlc.GetTreatmentHistoryByPatientID(r.Context(), patientID)
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"treatmentHistory": treatmentHistory}, nil)
}