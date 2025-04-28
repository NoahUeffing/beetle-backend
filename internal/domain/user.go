package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserAuthInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserCreateInput struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}

type User struct {
	Entity
	Email       string         `json:"email" validate:"email"`
	Password    string         `json:"-"`
	Username    string         `json:"username"`
	FirstName   NullableString `json:"first_name,omitempty"`
	LastName    NullableString `json:"last_name,omitempty"`
	Gender      NullableString `json:"gender,omitempty"`
	DateOfBirth *time.Time     `json:"date_of_birth,omitempty"`
	Country     NullableString `json:"country,omitempty"`
	City        NullableString `json:"city,omitempty"`
}

type UserAuthToken struct {
	Token string `json:"token"`
}

type IUserService interface {
	CreateUser(newUser *UserCreateInput) (*User, error)
	CheckPassword(user *User, password string) error
	CreateAuthToken(user *User) (*UserAuthToken, error)
	ReadByEmail(email string) (*User, error)
	ReadByID(id uuid.UUID) (*User, error)
	Update(user *User) (*User, error)
	// TODO: Add User Permissions
	// TODO: Add Reset Password
	// TODO: Add Delete User
}
