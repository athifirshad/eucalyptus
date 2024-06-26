package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/athifirshad/eucalyptus/db"
	"github.com/athifirshad/eucalyptus/internal/data"
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

func (app *application) GetAllDoctorInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Call the GetAllDoctorInfo function from your queries struct
	doctors, err := app.sqlc.GetAllDoctorInfo(r.Context())
	if err != nil {
		// Handle error, e.g., by sending an error response
		http.Error(w, "Error fetching doctors", http.StatusInternalServerError)
		return
	}

	// Marshal the doctors slice into JSON and write it to the response
	if err := json.NewEncoder(w).Encode(doctors); err != nil {
		// Handle error, e.g., by sending an error response
		http.Error(w, "Error marshalling doctors", http.StatusInternalServerError)
		return
	}
}

func (app *application) InsertAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	var params db.InsertAppointmentParams
	if err := app.readJSON(w, r, &params); err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error reading JSON:", err)
		return
	}

	// Insert the appointment using the InsertAppointment function from your queries struct
	err := app.sqlc.InsertAppointment(r.Context(), params)
	if err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error inserting appointment:", err)

		if errors.Is(err, pgx.ErrNoRows) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "appointment not inserted"}, nil)
			return
		}
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	// Respond with a success message
	app.writeJSON(w, http.StatusCreated, envelope{"message": "Appointment created successfully"}, nil)
}

func (app *application) getMedicationByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := app.contextGetUser(r)

	// Assuming the user object has a method or field to get the patientID
	patientID, err := app.models.Users.GetPatientIDByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Call the GetMedicationByPatientId function from your queries struct
	medications, err := app.sqlc.GetMedicationByPatientId(r.Context(), int32(patientID))
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"medications": medications}, nil)
}

func (app *application) getMedicalDirectivesByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := app.contextGetUser(r)

	// Assuming the user object has a method or field to get the patientID
	patientID, err := app.models.Users.GetPatientIDByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Call the GetMedicalDirectivesByPatientId function from your queries struct
	directives, err := app.sqlc.GetMedicalDirectivesByPatientId(r.Context(), int32(patientID))
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"directives": directives}, nil)
}

func (app *application) getVaccinationHistoryByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := app.contextGetUser(r)

	// Assuming the user object has a method or field to get the patientID
	patientID, err := app.models.Users.GetPatientIDByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Call the GetVaccinationHistoryByPatientId function from your queries struct
	history, err := app.sqlc.GetVaccinationHistoryByPatientId(r.Context(), int32(patientID))
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"history": history}, nil)
}

func (app *application) getAllergiesByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := app.contextGetUser(r)

	// Assuming the user object has a method or field to get the patientID
	patientID, err := app.models.Users.GetPatientIDByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Call the GetAllergiesByPatientId function from your queries struct
	allergies, err := app.sqlc.GetAllergiesByPatientId(r.Context(), int32(patientID))
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"allergies": allergies}, nil)
}

func (app *application) getFamilyMedicalHistoryByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := app.contextGetUser(r)

	// Assuming the user object has a method or field to get the patientID
	patientID, err := app.models.Users.GetPatientIDByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Call the GetFamilyMedicalHistoryByPatientId function from your queries struct
	familyHistory, err := app.sqlc.GetFamilyMedicalHistoryByPatientId(r.Context(), int32(patientID))
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"family_history": familyHistory}, nil)
}

func (app *application) getSocialHistoryByPatientIdHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := app.contextGetUser(r)

	// Assuming the user object has a method or field to get the patientID
	patientID, err := app.models.Users.GetPatientIDByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Call the GetSocialHistoryByPatientId function from your queries struct
	socialHistory, err := app.sqlc.GetSocialHistoryByPatientId(r.Context(), int32(patientID))
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"social_history": socialHistory}, nil)
}

func (app *application) GetTreatmentHistoryByPatientIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := app.contextGetUser(r)

	// Assuming the user object has a method or field to get the patientID
	patientID, err := app.models.Users.GetPatientIDByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Call the GetTreatmentHistoryByPatientID function from your queries struct
	treatmentHistory, err := app.sqlc.GetTreatmentHistoryByPatientID(r.Context(), int32(patientID))
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "internal server error"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"treatmentHistory": treatmentHistory}, nil)
}

func (app *application) registerDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"usertype"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		UserType:  data.UserType(input.UserType),
		Activated: false,
	}

	err = user.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//	userType, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = app.models.Users.CreateUser(user, "doctor")

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := map[string]any{
		"activationToken": token.Plaintext,
		"userID":          user.ID,
	}

	err = app.mailer.Send(user.Email, "user_welcome.htm", data)
	if err != nil {
		app.logger.Sugar().Error(err, nil)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) DoctorProfileHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := app.contextGetUser(r)

	var input db.CreateUserProfileParams
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Use the UserID from the current user.
	input.UserID = int32(currentUser.ID)

	profile, err := app.sqlc.GetProfileByUserId(r.Context(), input.UserID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
		return
	}

	err = app.sqlc.CreateUserProfile(r.Context(), input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"message": "User profile created successfully"}, nil)
}
