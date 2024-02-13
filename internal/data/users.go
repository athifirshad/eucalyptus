package data

import (
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserTypeEnum string

const (
	UserTypeEnumPatient       UserTypeEnum = "patient"
	UserTypeEnumDoctor        UserTypeEnum = "doctor"
	UserTypeEnumAdministrator UserTypeEnum = "administrator"
)

type NullUserTypeEnum struct {
	UserTypeEnum UserTypeEnum
	Valid        bool
}

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
