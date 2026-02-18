package internal

import (
	"fmt"
	"strings"
	"time"
)

type UserRepository struct {
	users map[string]*User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*User),
	}
}

func (r *UserRepository) registerUser(user UserRegister) (User, error) {
	userAppend := &User{
		ID:          1,
		FullName:    user.FullName,
		Username:    strings.Split(user.Email, "@")[0],
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	key := fmt.Sprintf("id-%s", user.Email)
	r.users[key] = userAppend

	return *r.users[key], nil
}
