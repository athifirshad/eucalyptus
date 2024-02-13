// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getDoctorById = `-- name: GetDoctorById :one
SELECT doctor_id, profile_id, specialization, hospital_id, available_consultation_time FROM doctor WHERE doctor_id = $1
`

func (q *Queries) GetDoctorById(ctx context.Context, doctorID int32) (Doctor, error) {
	row := q.db.QueryRow(ctx, getDoctorById, doctorID)
	var i Doctor
	err := row.Scan(
		&i.DoctorID,
		&i.ProfileID,
		&i.Specialization,
		&i.HospitalID,
		&i.AvailableConsultationTime,
	)
	return i, err
}

const getHealthRecordByRecordId = `-- name: GetHealthRecordByRecordId :one
SELECT record_id, patient_id, weight, height, treatment_history, medical_directives, vaccination_history, allergies, family_medical_history, social_history, created_at, updated_at FROM health_record WHERE record_id = $1
`

func (q *Queries) GetHealthRecordByRecordId(ctx context.Context, recordID int32) (HealthRecord, error) {
	row := q.db.QueryRow(ctx, getHealthRecordByRecordId, recordID)
	var i HealthRecord
	err := row.Scan(
		&i.RecordID,
		&i.PatientID,
		&i.Weight,
		&i.Height,
		&i.TreatmentHistory,
		&i.MedicalDirectives,
		&i.VaccinationHistory,
		&i.Allergies,
		&i.FamilyMedicalHistory,
		&i.SocialHistory,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getHealthRecordsByPatientId = `-- name: GetHealthRecordsByPatientId :many
SELECT record_id, patient_id, weight, height, treatment_history, medical_directives, vaccination_history, allergies, family_medical_history, social_history, created_at, updated_at FROM health_record WHERE patient_id = $1
`

func (q *Queries) GetHealthRecordsByPatientId(ctx context.Context, patientID pgtype.Int4) ([]HealthRecord, error) {
	rows, err := q.db.Query(ctx, getHealthRecordsByPatientId, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []HealthRecord
	for rows.Next() {
		var i HealthRecord
		if err := rows.Scan(
			&i.RecordID,
			&i.PatientID,
			&i.Weight,
			&i.Height,
			&i.TreatmentHistory,
			&i.MedicalDirectives,
			&i.VaccinationHistory,
			&i.Allergies,
			&i.FamilyMedicalHistory,
			&i.SocialHistory,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getHospitalByHospitalId = `-- name: GetHospitalByHospitalId :one
SELECT hospital_id, hospital_name, address FROM hospital WHERE hospital_id = $1
`

func (q *Queries) GetHospitalByHospitalId(ctx context.Context, hospitalID int32) (Hospital, error) {
	row := q.db.QueryRow(ctx, getHospitalByHospitalId, hospitalID)
	var i Hospital
	err := row.Scan(&i.HospitalID, &i.HospitalName, &i.Address)
	return i, err
}

const getMedicationsByPrescriptionId = `-- name: GetMedicationsByPrescriptionId :many
SELECT medication_id, prescription_id, medication_name, dosage, frequency, start_date, end_date, instructions FROM medication WHERE prescription_id = $1
`

func (q *Queries) GetMedicationsByPrescriptionId(ctx context.Context, prescriptionID pgtype.Int4) ([]Medication, error) {
	rows, err := q.db.Query(ctx, getMedicationsByPrescriptionId, prescriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Medication
	for rows.Next() {
		var i Medication
		if err := rows.Scan(
			&i.MedicationID,
			&i.PrescriptionID,
			&i.MedicationName,
			&i.Dosage,
			&i.Frequency,
			&i.StartDate,
			&i.EndDate,
			&i.Instructions,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPrescriptionsByPatientId = `-- name: GetPrescriptionsByPatientId :many
SELECT prescription_id, doctor_id, patient_id, diagnosis FROM prescription WHERE patient_id = $1
`

func (q *Queries) GetPrescriptionsByPatientId(ctx context.Context, patientID pgtype.Int4) ([]Prescription, error) {
	rows, err := q.db.Query(ctx, getPrescriptionsByPatientId, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Prescription
	for rows.Next() {
		var i Prescription
		if err := rows.Scan(
			&i.PrescriptionID,
			&i.DoctorID,
			&i.PatientID,
			&i.Diagnosis,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProfileByUserId = `-- name: GetProfileByUserId :one
SELECT profile_id, user_id, name, date_of_birth, gender, address, phone_number, email, marital_status, nationality, language_preference FROM profile WHERE user_id = $1
`

func (q *Queries) GetProfileByUserId(ctx context.Context, userID pgtype.Int4) (Profile, error) {
	row := q.db.QueryRow(ctx, getProfileByUserId, userID)
	var i Profile
	err := row.Scan(
		&i.ProfileID,
		&i.UserID,
		&i.Name,
		&i.DateOfBirth,
		&i.Gender,
		&i.Address,
		&i.PhoneNumber,
		&i.Email,
		&i.MaritalStatus,
		&i.Nationality,
		&i.LanguagePreference,
	)
	return i, err
}

const getUsersByUsername = `-- name: GetUsersByUsername :one
SELECT user_id, username, password, user_type FROM app_user WHERE user_id = $1
`

func (q *Queries) GetUsersByUsername(ctx context.Context, userID int32) (AppUser, error) {
	row := q.db.QueryRow(ctx, getUsersByUsername, userID)
	var i AppUser
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Password,
		&i.UserType,
	)
	return i, err
}

const insertAppointment = `-- name: InsertAppointment :exec
INSERT INTO appointment (doctor_id, patient_id, appointment_date, status)
VALUES ($1, $2, $3, $4)
RETURNING appointment_id
`

type InsertAppointmentParams struct {
	DoctorID        pgtype.Int4
	PatientID       pgtype.Int4
	AppointmentDate pgtype.Timestamp
	Status          NullAppointmentStatus
}

func (q *Queries) InsertAppointment(ctx context.Context, arg InsertAppointmentParams) error {
	_, err := q.db.Exec(ctx, insertAppointment,
		arg.DoctorID,
		arg.PatientID,
		arg.AppointmentDate,
		arg.Status,
	)
	return err
}

const insertHealthRecord = `-- name: InsertHealthRecord :exec
INSERT INTO health_record (patient_id, weight, height, treatment_history, medical_directives, vaccination_history, allergies, family_medical_history, social_history)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING record_id
`

type InsertHealthRecordParams struct {
	PatientID            pgtype.Int4
	Weight               pgtype.Numeric
	Height               pgtype.Numeric
	TreatmentHistory     pgtype.Text
	MedicalDirectives    pgtype.Text
	VaccinationHistory   pgtype.Text
	Allergies            pgtype.Text
	FamilyMedicalHistory pgtype.Text
	SocialHistory        pgtype.Text
}

func (q *Queries) InsertHealthRecord(ctx context.Context, arg InsertHealthRecordParams) error {
	_, err := q.db.Exec(ctx, insertHealthRecord,
		arg.PatientID,
		arg.Weight,
		arg.Height,
		arg.TreatmentHistory,
		arg.MedicalDirectives,
		arg.VaccinationHistory,
		arg.Allergies,
		arg.FamilyMedicalHistory,
		arg.SocialHistory,
	)
	return err
}

const insertMedication = `-- name: InsertMedication :exec
INSERT INTO medication (prescription_id, medication_name, dosage, frequency, start_date, end_date, instructions)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING medication_id
`

type InsertMedicationParams struct {
	PrescriptionID pgtype.Int4
	MedicationName pgtype.Text
	Dosage         pgtype.Text
	Frequency      pgtype.Text
	StartDate      pgtype.Date
	EndDate        pgtype.Date
	Instructions   pgtype.Text
}

func (q *Queries) InsertMedication(ctx context.Context, arg InsertMedicationParams) error {
	_, err := q.db.Exec(ctx, insertMedication,
		arg.PrescriptionID,
		arg.MedicationName,
		arg.Dosage,
		arg.Frequency,
		arg.StartDate,
		arg.EndDate,
		arg.Instructions,
	)
	return err
}

const insertPrescription = `-- name: InsertPrescription :exec
INSERT INTO prescription (doctor_id, patient_id, diagnosis)
VALUES ($1, $2, $3)
RETURNING prescription_id
`

type InsertPrescriptionParams struct {
	DoctorID  pgtype.Int4
	PatientID pgtype.Int4
	Diagnosis pgtype.Text
}

func (q *Queries) InsertPrescription(ctx context.Context, arg InsertPrescriptionParams) error {
	_, err := q.db.Exec(ctx, insertPrescription, arg.DoctorID, arg.PatientID, arg.Diagnosis)
	return err
}

const updateAppointmentStatus = `-- name: UpdateAppointmentStatus :exec
UPDATE appointment
SET status = $2
WHERE appointment_id = $1
`

type UpdateAppointmentStatusParams struct {
	AppointmentID int32
	Status        NullAppointmentStatus
}

func (q *Queries) UpdateAppointmentStatus(ctx context.Context, arg UpdateAppointmentStatusParams) error {
	_, err := q.db.Exec(ctx, updateAppointmentStatus, arg.AppointmentID, arg.Status)
	return err
}

const updateHealthRecord = `-- name: UpdateHealthRecord :exec
UPDATE health_record
SET weight = $2, height = $3, treatment_history = $4, medical_directives = $5, vaccination_history = $6, allergies = $7, family_medical_history = $8, social_history = $9
WHERE record_id = $1
`

type UpdateHealthRecordParams struct {
	RecordID             int32
	Weight               pgtype.Numeric
	Height               pgtype.Numeric
	TreatmentHistory     pgtype.Text
	MedicalDirectives    pgtype.Text
	VaccinationHistory   pgtype.Text
	Allergies            pgtype.Text
	FamilyMedicalHistory pgtype.Text
	SocialHistory        pgtype.Text
}

func (q *Queries) UpdateHealthRecord(ctx context.Context, arg UpdateHealthRecordParams) error {
	_, err := q.db.Exec(ctx, updateHealthRecord,
		arg.RecordID,
		arg.Weight,
		arg.Height,
		arg.TreatmentHistory,
		arg.MedicalDirectives,
		arg.VaccinationHistory,
		arg.Allergies,
		arg.FamilyMedicalHistory,
		arg.SocialHistory,
	)
	return err
}

const updateMedication = `-- name: UpdateMedication :exec
UPDATE medication
SET medication_name = $2, dosage = $3, frequency = $4, start_date = $5, end_date = $6, instructions = $7
WHERE medication_id = $1
`

type UpdateMedicationParams struct {
	MedicationID   int32
	MedicationName pgtype.Text
	Dosage         pgtype.Text
	Frequency      pgtype.Text
	StartDate      pgtype.Date
	EndDate        pgtype.Date
	Instructions   pgtype.Text
}

func (q *Queries) UpdateMedication(ctx context.Context, arg UpdateMedicationParams) error {
	_, err := q.db.Exec(ctx, updateMedication,
		arg.MedicationID,
		arg.MedicationName,
		arg.Dosage,
		arg.Frequency,
		arg.StartDate,
		arg.EndDate,
		arg.Instructions,
	)
	return err
}

const updatePrescription = `-- name: UpdatePrescription :exec
UPDATE prescription
SET diagnosis = $2
WHERE prescription_id = $1
`

type UpdatePrescriptionParams struct {
	PrescriptionID int32
	Diagnosis      pgtype.Text
}

func (q *Queries) UpdatePrescription(ctx context.Context, arg UpdatePrescriptionParams) error {
	_, err := q.db.Exec(ctx, updatePrescription, arg.PrescriptionID, arg.Diagnosis)
	return err
}
