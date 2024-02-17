// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCanceled  AppointmentStatus = "canceled"
)

func (e *AppointmentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AppointmentStatus(s)
	case string:
		*e = AppointmentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for AppointmentStatus: %T", src)
	}
	return nil
}

type NullAppointmentStatus struct {
	AppointmentStatus AppointmentStatus `json:"appointment_status"`
	Valid             bool              `json:"valid"` // Valid is true if AppointmentStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAppointmentStatus) Scan(value interface{}) error {
	if value == nil {
		ns.AppointmentStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AppointmentStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAppointmentStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AppointmentStatus), nil
}

type UserTypeEnum string

const (
	UserTypeEnumPatient       UserTypeEnum = "patient"
	UserTypeEnumDoctor        UserTypeEnum = "doctor"
	UserTypeEnumAdministrator UserTypeEnum = "administrator"
)

func (e *UserTypeEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserTypeEnum(s)
	case string:
		*e = UserTypeEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for UserTypeEnum: %T", src)
	}
	return nil
}

type NullUserTypeEnum struct {
	UserTypeEnum UserTypeEnum `json:"user_type_enum"`
	Valid        bool         `json:"valid"` // Valid is true if UserTypeEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserTypeEnum) Scan(value interface{}) error {
	if value == nil {
		ns.UserTypeEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserTypeEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserTypeEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserTypeEnum), nil
}

type AppUser struct {
	UserID   int32            `json:"user_id"`
	Username pgtype.Text      `json:"username"`
	Password pgtype.Text      `json:"password"`
	UserType NullUserTypeEnum `json:"user_type"`
}

type Appointment struct {
	AppointmentID   int32                 `json:"appointment_id"`
	DoctorID        pgtype.Int4           `json:"doctor_id"`
	PatientID       pgtype.Int4           `json:"patient_id"`
	AppointmentDate pgtype.Timestamp      `json:"appointment_date"`
	Status          NullAppointmentStatus `json:"status"`
}

type Doctor struct {
	DoctorID                  int32       `json:"doctor_id"`
	ProfileID                 pgtype.Int4 `json:"profile_id"`
	Specialization            pgtype.Text `json:"specialization"`
	HospitalID                pgtype.Int4 `json:"hospital_id"`
	AvailableConsultationTime pgtype.Text `json:"available_consultation_time"`
}

type HealthRecord struct {
	RecordID             int32            `json:"record_id"`
	PatientID            pgtype.Int4      `json:"patient_id"`
	Weight               pgtype.Numeric   `json:"weight"`
	Height               pgtype.Numeric   `json:"height"`
	TreatmentHistory     pgtype.Text      `json:"treatment_history"`
	MedicalDirectives    pgtype.Text      `json:"medical_directives"`
	VaccinationHistory   pgtype.Text      `json:"vaccination_history"`
	Allergies            pgtype.Text      `json:"allergies"`
	FamilyMedicalHistory pgtype.Text      `json:"family_medical_history"`
	SocialHistory        pgtype.Text      `json:"social_history"`
	CreatedAt            pgtype.Timestamp `json:"created_at"`
	UpdatedAt            pgtype.Timestamp `json:"updated_at"`
}

type Hospital struct {
	HospitalID   int32       `json:"hospital_id"`
	HospitalName pgtype.Text `json:"hospital_name"`
	Address      pgtype.Text `json:"address"`
}

type Medication struct {
	MedicationID   int32       `json:"medication_id"`
	PrescriptionID pgtype.Int4 `json:"prescription_id"`
	MedicationName pgtype.Text `json:"medication_name"`
	Dosage         pgtype.Text `json:"dosage"`
	Frequency      pgtype.Text `json:"frequency"`
	StartDate      pgtype.Date `json:"start_date"`
	EndDate        pgtype.Date `json:"end_date"`
	Instructions   pgtype.Text `json:"instructions"`
}

type Patient struct {
	PatientID int32       `json:"patient_id"`
	ProfileID pgtype.Int4 `json:"profile_id"`
}

type Prescription struct {
	PrescriptionID int32       `json:"prescription_id"`
	DoctorID       pgtype.Int4 `json:"doctor_id"`
	PatientID      pgtype.Int4 `json:"patient_id"`
	Diagnosis      pgtype.Text `json:"diagnosis"`
}

type Profile struct {
	ProfileID          int32       `json:"profile_id"`
	UserID             pgtype.Int4 `json:"user_id"`
	Name               pgtype.Text `json:"name"`
	DateOfBirth        pgtype.Date `json:"date_of_birth"`
	Gender             pgtype.Text `json:"gender"`
	Address            pgtype.Text `json:"address"`
	PhoneNumber        pgtype.Text `json:"phone_number"`
	Email              pgtype.Text `json:"email"`
	MaritalStatus      pgtype.Text `json:"marital_status"`
	Nationality        pgtype.Text `json:"nationality"`
	LanguagePreference pgtype.Text `json:"language_preference"`
}
