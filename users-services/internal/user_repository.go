package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func hashingPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (r *UserRepository) registerUser(c context.Context, user UserRegister) (*User, error) {
	query := make(map[string]string)
	var username, hashing string
	var userId int32

	if user.Password != nil {
		hashing = hashingPassword(*user.Password)
	}

	query["query_insert"] = `
		insert into users (
		    full_name,
			username,
		    email,
		    password,
		    phone_number,
		    is_active
		) values (
			$1, $2, $3, $4, $5, $6
		) RETURNING id,username
    `

	err := r.db.QueryRowContext(c, query["query_insert"], user.FullName, strings.Split(user.Email, "@")[0], user.Email, hashing, user.PhoneNumber, true).Scan(&userId, &username)

	if err != nil {
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == "23505" {
			switch pgErr.Constraint {
			case "users_email_key", "idx_users_email":
				return nil, fmt.Errorf("email already registered")
			case "idx_users_phone_number":
				return nil, fmt.Errorf("phone number already registered")
			}
		}
		return nil, fmt.Errorf("could not insert user: %v", err)
	}
	return &User{
		ID:          userId,
		FullName:    user.FullName,
		Username:    username,
		Email:       user.Email,
		Password:    nil,
		PhoneNumber: user.PhoneNumber,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
func (r *UserRepository) loginUser(c context.Context, email string, password string) (string, error) {
	query := `
        SELECT id, password FROM users WHERE email = $1 AND is_deleted = false AND is_active = true;
    `

	var userId int32
	var hashedPassword string

	err := r.db.QueryRowContext(c, query, email).Scan(&userId, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("invalid email or password")
		}
		return "", fmt.Errorf("could not query user: %v", err)
	}

	if err := verifyPassword(hashedPassword, password); err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	token, err := GenerateToken(fmt.Sprintf("%d", userId))
	if err != nil {
		return "", fmt.Errorf("could not generate token: %v", err)
	}

	return token, nil
}

func (r *UserRepository) findMe(c context.Context, userID int32) (*User, error) {
	var user User
	query := `
		SELECT 
		    	full_name,
		    	email,
		    	phone_number,
		    	is_active,
		    	created_at 
		FROM users WHERE id = $1
		AND is_deleted = false;
    `
	err := r.db.QueryRowContext(c, query, userID).Scan(&user.FullName, &user.Email, &user.PhoneNumber, &user.IsActive, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("could not query user: %v", err)
	}
	return &user, nil
}
