-- name: GetUsersByUsername :one
SELECT * FROM app_user WHERE user_id = $1;

-- name: GetProfileByUserId :one
SELECT * FROM profile WHERE user_id = $1;

-- name: GetHealthRecordsByPatientId :many
SELECT * FROM health_record WHERE patient_id = $1;

-- name: GetHealthRecordByRecordId :one
SELECT * FROM health_record WHERE record_id = $1;

-- name: GetHospitalByHospitalId :one
SELECT * FROM hospital WHERE hospital_id = $1;

-- name: GetDoctorById :one
SELECT * FROM doctor WHERE doctor_id = $1;

-- name: GetPrescriptionsByPatientId :many
SELECT * FROM prescription WHERE patient_id = $1;

-- name: GetMedicationsByPrescriptionId :many
SELECT * FROM medication WHERE prescription_id = $1;


-- name: InsertAppointment :exec
INSERT INTO appointment (doctor_id, patient_id, appointment_date, status)
VALUES ($1, $2, $3, $4)
RETURNING appointment_id;

-- name: UpdateAppointmentStatus :exec
UPDATE appointment
SET status = $2
WHERE appointment_id = $1;

-- name: InsertPrescription :exec
INSERT INTO prescription (doctor_id, patient_id, diagnosis)
VALUES ($1, $2, $3)
RETURNING prescription_id;

-- name: UpdatePrescription :exec
UPDATE prescription
SET diagnosis = $2
WHERE prescription_id = $1;

-- name: InsertMedication :exec
INSERT INTO medication (prescription_id, medication_name, dosage, frequency, start_date, end_date, instructions)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING medication_id;

-- name: UpdateMedication :exec
UPDATE medication
SET medication_name = $2, dosage = $3, frequency = $4, start_date = $5, end_date = $6, instructions = $7
WHERE medication_id = $1;


-- name: InsertHealthRecord :exec
INSERT INTO health_record (patient_id, weight, height, treatment_history, medical_directives, vaccination_history, allergies, family_medical_history, social_history)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING record_id;

-- name: UpdateHealthRecord :exec
UPDATE health_record
SET weight = $2, height = $3, treatment_history = $4, medical_directives = $5, vaccination_history = $6, allergies = $7, family_medical_history = $8, social_history = $9
WHERE record_id = $1;
