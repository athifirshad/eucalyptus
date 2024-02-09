package data

import "time"

// PatientHealthRecord represents the complete health record of a patient,
// including personal information, medical history, and treatment plans.
type PatientHealthRecord struct {
	ID                   int64     `json:"id"`
	TreatmentHistory     string    `json:"treatment_history"`
	MedicalDirectives    string    `json:"medical_directives"`
	VaccinationHistory   string    `json:"vaccination_history"`
	Allergies            string    `json:"allergies"`
	FamilyMedicalHistory string    `json:"family_medical_history"`
	SocialHistory        string    `json:"social_history"`
	ReviewOfSystems      string    `json:"review_of_systems"`
	PhysicalExaminations string    `json:"physical_examinations"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
