package main

import (
	"net/http"
	"time"

	"github.com/athifirshad/eucalyptus/internal/data"
)

func (app *application) showMedicalRecord(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a mock patient health record 
	mockRecord := data.PatientHealthRecord{
		ID:                   102,
		TreatmentHistory:     "Initial check-up, Annual flu shot",
		MedicalDirectives:    "Do not smoke",
		VaccinationHistory:   "Flu shot  2023",
		Allergies:            "Peanuts",
		FamilyMedicalHistory: "Grandfather had heart disease",
		SocialHistory:        "Regular exercise",
		ReviewOfSystems:      "No significant findings",
		PhysicalExaminations: "Normal vital signs",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}


	if id != mockRecord.ID {
		http.NotFound(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"mockrecord":mockRecord}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
