package data

import (
	"context"
	"errors"

	"github.com/jackc/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error{
	hash, err:= bcrypt.GenerateFromPassword([]byte(plaintextPassword),12)
	if err!= nil{
		return err
	}

	p.plaintext=&plaintextPassword
	p.hash=hash

	return nil
}


func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}


//HACKY SQLC GENERATED CODE BELOW
type User struct {
	UserID       int64              `json:"user_id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	PasswordHash []byte             `json:"-"`
	Activated    bool               `json:"activated"`
	Version      int32              `json:"-"`
	UserType     NullUserTypeEnum   `json:"user_type"`
}

const createAdminUser = `-- name: CreateAdminUser :exec
INSERT INTO users (name, email, password_hash, user_type)
VALUES ($1, $2, $3, 'administrator')
`

type CreateAdminUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"-"`
}

func (q *Queries) CreateAdminUser(ctx context.Context, arg CreateAdminUserParams) error {
	_, err := q.db.Exec(ctx, createAdminUser, arg.Name, arg.Email, arg.PasswordHash)
	return err
}

const createDoctorUser = `-- name: CreateDoctorUser :exec
INSERT INTO users (name, email, password_hash, user_type)
VALUES ($1, $2, $3, 'doctor')
`

type CreateDoctorUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"-"`
}

func (q *Queries) CreateDoctorUser(ctx context.Context, arg CreateDoctorUserParams) error {
	_, err := q.db.Exec(ctx, createDoctorUser, arg.Name, arg.Email, arg.PasswordHash)
	return err
}

const createPatientUser = `-- name: CreatePatientUser :exec
INSERT INTO users (name, email, password_hash, user_type)
VALUES ($1, $2, $3, 'patient')
`

type CreatePatientUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"-"`
}

func (q *Queries) CreatePatientUser(ctx context.Context, arg CreatePatientUserParams) error {
	_, err := q.db.Exec(ctx, createPatientUser, arg.Name, arg.Email, arg.PasswordHash)
	return err
}
