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

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

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

type UserModel struct {
	*User
}

func (um *UserModel) SetPassword(plaintextPassword string) error {
	pwd := &password{}
	err := pwd.Set(plaintextPassword)
	if err != nil {
		return err
	}
	um.PasswordHash = pwd.hash
	return nil
}


var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// HACKY SQLC GENERATED CODE BELOW
type User struct {
	UserID       int64              `json:"user_id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	PasswordHash []byte             `json:"-"`
	Activated    bool               `json:"activated"`
	Version      int32              `json:"-"`
	UserType     string             `json:"user_type"` //enums are a bitch to work with
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

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET name = $1, email = $2, password_hash = $3, activated = $4, user_type=$5, version = version + 1
WHERE user_id = $6 AND version = $7
RETURNING version
`

type UpdateUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"-"`
	Activated    bool   `json:"activated"`
	UserType     string `json:"user_type"` //fucking enums 
	UserID       int64  `json:"user_id"`
	Version      int32  `json:"version"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
		arg.Activated,
		arg.UserType,
		arg.UserID,
		arg.Version,
	)
	return err
}

const getByEmail = `-- name: GetByEmail :one
SELECT user_id, created_at, name, email, password_hash, activated, version, user_type
FROM users
WHERE EMAIL = $1
`

func (q *Queries) GetByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
		&i.UserType,
	)
	return i, err
}
