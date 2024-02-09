// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

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
	UserTypeEnum UserTypeEnum
	Valid        bool // Valid is true if UserTypeEnum is not NULL
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

type Doctor struct {
	DoctorID  int32
	ProfileID pgtype.Int4
}

type HealthRecord struct {
	RecordID             int32
	PatientID            pgtype.Int4
	TreatmentHistory     pgtype.Text
	MedicalDirectives    pgtype.Text
	VaccinationHistory   pgtype.Text
	Allergies            pgtype.Text
	FamilyMedicalHistory pgtype.Text
	SocialHistory        pgtype.Text
	ReviewOfSystems      pgtype.Text
	PhysicalExaminations pgtype.Text
	CreatedAt            pgtype.Timestamp
	UpdatedAt            pgtype.Timestamp
}

type Medication struct {
	MedicationID   int32
	PrescriptionID pgtype.Int4
	MedicationName pgtype.Text
	Dosage         pgtype.Text
	Frequency      pgtype.Text
	StartDate      pgtype.Date
	EndDate        pgtype.Date
	Instructions   pgtype.Text
}

type Patient struct {
	PatientID int32
	ProfileID pgtype.Int4
}

type Prescription struct {
	PrescriptionID int32
	DoctorID       pgtype.Int4
	PatientID      pgtype.Int4
	Diagnosis      pgtype.Text
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
}

type Profile struct {
	ProfileID          int32
	UserID             pgtype.Int4
	Name               pgtype.Text
	DateOfBirth        pgtype.Date
	Gender             pgtype.Text
	Address            pgtype.Text
	PhoneNumber        pgtype.Text
	Email              pgtype.Text
	MaritalStatus      pgtype.Text
	Nationality        pgtype.Text
	LanguagePreference pgtype.Text
}

type User struct {
	UserID   int32
	Username pgtype.Text
	Password pgtype.Text
	UserType NullUserTypeEnum
}
