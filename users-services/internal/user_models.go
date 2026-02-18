package internal

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID          int32     `json:"userId"`
	FullName    string    `json:"fullName"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    *string   `json:"password,omitempty"`
	PhoneNumber string    `json:"phoneNumber"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserRegister struct {
	FullName    string  `json:"fullName"`
	Email       string  `json:"email"`
	Password    *string `json:"password,omitempty"`
	PhoneNumber string  `json:"phoneNumber"`
}
type Repository interface {
	registerUser(user UserRegister) (User, error)
}

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[0-9]{8,15}$`)
)

func (u *UserRegister) Validate() error {
	if err := validateFullName(u.FullName); err != nil {
		return err
	}
	if err := validateEmail(u.Email); err != nil {
		return err
	}
	if err := validatePhone(u.PhoneNumber); err != nil {
		return err
	}
	return nil
}

func validateFullName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 4 {
		return errors.New("full name minimal 4 characters")
	}
	if len(name) > 100 {
		return errors.New("full name max 100 characters")
	}
	return nil
}

func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("format email not valid")
	}
	return nil
}

func validatePhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return errors.New("format phone number not valid")
	}
	return nil
}
