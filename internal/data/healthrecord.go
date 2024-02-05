package data

import "time"

type PatientHealthRecord struct {
	ID                   int64     `json:"id"`
	Name                 string    `json:"name"`
	DateOfBirth          time.Time `json:"date_of_birth"`
	Gender               string    `json:"gender"`
	Address              string    `json:"address"`
	PhoneNumber          string    `json:"phone_number"`
	Email                string    `json:"email"`
	MaritalStatus        string    `json:"marital_status"`
	Nationality          string    `json:"nationality"`
	LanguagePreference   string    `json:"language_preference"`
	TreatmentHistory     []string  `json:"treatment_history"`
	MedicalDirectives    []string  `json:"medical_directives"`
	VaccinationHistory   []string  `json:"vaccination_history"`
	Allergies            []string  `json:"allergies"`
	FamilyMedicalHistory []string  `json:"family_medical_history"`
	SocialHistory        []string  `json:"social_history"`
	ReviewOfSystems      []string  `json:"review_of_systems"`
	PhysicalExaminations []string  `json:"physical_examinations"`
	Diagnoses            []string  `json:"diagnoses"`
	Procedures           []string  `json:"procedures"`
	PlansAndOrders       []string  `json:"plans_and_orders"`
	Notes                []string  `json:"notes"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
