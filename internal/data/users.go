package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"
  //  "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	dbPool *pgxpool.Pool
}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
	UserType  UserType  `json:"usertype"`
}

type UserType string

type password struct {
	plaintext *string
	hash      []byte
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("no record found")
	ErrEditConflict   = errors.New("unable to edit due to conflict")
)

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

func (m UserModel) CreateUser(user *User, usertype string) error {
	query := `
	INSERT INTO users (name, email, password_hash, activated, user_type)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING user_id, created_at, version, user_type`
	args := []any{user.Name, user.Email, user.Password.hash, user.Activated, usertype}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.dbPool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version, &user.UserType)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT user_id, created_at, name, email, password_hash, activated, version, user_type
	FROM users
	WHERE email = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.dbPool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.UserType,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
	WHERE user_id = $5 AND version = $6
	RETURNING version`
	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
		//user.UserType,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.dbPool.QueryRow(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	if tokenPlaintext != "" {
		fmt.Printf("Token Plaintext: %x\n", tokenPlaintext)
	} else {
		fmt.Println("Token Plaintext is empty")
	}
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	fmt.Printf("Token Hash: %x\n", tokenHash[:])
	// Set up the SQL query.
	query := `
	SELECT users.user_id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version, users.user_type
	FROM users
	INNER JOIN tokens
	ON users.user_id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`
	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.dbPool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.UserType,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.
	return &user, nil
}
