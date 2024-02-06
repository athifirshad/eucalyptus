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
		Name:                 "John Doe",
		DateOfBirth:          time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:               "Male",
		Address:              "123 Main St",
		PhoneNumber:          "+1234567890",
		Email:                "john.doe@example.com",
		MaritalStatus:        "Single",
		Nationality:          "American",
		LanguagePreference:   "English",
		TreatmentHistory:     []string{"Initial check-up", "Annual flu shot"},
		MedicalDirectives:    []string{"Do not smoke"},
		VaccinationHistory:   []string{"Flu shot 2023"},
		Allergies:            []string{"Peanuts"},
		FamilyMedicalHistory: []string{"Grandfather had heart disease"},
		SocialHistory:        []string{"Regular exercise"},
		ReviewOfSystems:      []string{"No significant findings"},
		PhysicalExaminations: []string{"Normal vital signs"},
		Diagnoses:            []string{"No recent diagnoses"},
		Procedures:           []string{"No recent procedures"},
		PlansAndOrders:       []string{"Routine check-ups every six months"},
		Notes:                []string{"Patient is generally healthy"},
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
