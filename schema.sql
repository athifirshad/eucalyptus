-- Define the user_type enum type
CREATE TYPE user_type_enum AS ENUM ('patient', 'doctor', 'administrator');

-- Create the user table
CREATE TABLE "user" (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR,
    password VARCHAR,
    user_type user_type_enum
);

CREATE TABLE profile (
    profile_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user"(user_id),
    name VARCHAR,
    date_of_birth DATE,
    gender VARCHAR,
    address VARCHAR,
    phone_number VARCHAR,
    email VARCHAR,
    marital_status VARCHAR,
    nationality VARCHAR,
    language_preference VARCHAR
);

CREATE TABLE patient (
    patient_id SERIAL PRIMARY KEY,
    profile_id INT REFERENCES profile(profile_id)
);

CREATE TABLE health_record (
    record_id SERIAL PRIMARY KEY,
    patient_id INT REFERENCES patient(patient_id),
    treatment_history TEXT,
    medical_directives TEXT,
    vaccination_history TEXT,
    allergies TEXT,
    family_medical_history TEXT,
    social_history TEXT,
    review_of_systems TEXT,
    physical_examinations TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE doctor (
    doctor_id SERIAL PRIMARY KEY,
    profile_id INT REFERENCES profile(profile_id)
);

CREATE TABLE prescription (
    prescription_id SERIAL PRIMARY KEY,
    doctor_id INT REFERENCES doctor(doctor_id),
    patient_id INT REFERENCES patient(patient_id),
    diagnosis TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE medication (
    medication_id SERIAL PRIMARY KEY,
    prescription_id INT REFERENCES prescription(prescription_id),
    medication_name VARCHAR,
    dosage VARCHAR,
    frequency VARCHAR,
    start_date DATE,
    end_date DATE,
    instructions TEXT
);
