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

-- name: GetTreatmentHistoryByPatientID :many
SELECT * FROM treatment_history WHERE patient_id = $1;

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

-- name: GetAllDoctorInfo :many
select profile.name, doctor.doctor_id, doctor.profile_id ,doctor.specialization,doctor.hospital_id, doctor.available_consultation_time, hospital.hospital_name,hospital.address from doctor 
inner join hospital on doctor.hospital_id=hospital.hospital_id
inner join profile on doctor.profile_id=profile.profile_id;

-- name: GetMedicationByPatientId :many
SELECT *
FROM medication
INNER JOIN prescription ON medication.prescription_id = prescription.prescription_id
INNER JOIN patient ON prescription.patient_id = patient.patient_id
WHERE patient.patient_id = $1;

-- name: GetMedicalDirectivesByPatientId :many
SELECT * FROM medical_directives WHERE patient_id = $1;

-- name: GetVaccinationHistoryByPatientId :many
SELECT * FROM vaccination_history WHERE patient_id = $1;

-- name: GetAllergiesByPatientId :many
SELECT * FROM allergies WHERE patient_id = $1;

-- name: GetFamilyMedicalHistoryByPatientId :many
SELECT * FROM family_medical_history WHERE patient_id = $1;

-- name: GetSocialHistoryByPatientId :many
SELECT * FROM social_history WHERE patient_id = $1;


-- name: CreateUserProfile :exec
INSERT INTO profile (user_id, name, date_of_birth, gender, address, phone_number, email, marital_status, nationality, language_preference)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);


-- name: FindPrescriptions :many
SELECT p.prescription_id, p.diagnosis, m.medication_name, m.dosage, m.frequency, m.start_date, m.end_date, m.instructions
FROM prescription p
JOIN patient pt ON p.patient_id = pt.patient_id
JOIN profile pr ON pt.profile_id = pr.profile_id
JOIN medication m ON p.prescription_id = m.prescription_id
WHERE pr.user_id = $1;


















