package domain

import (
	"time"
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
	Email       string     `json:"email" validate:"email"`
	Password    string     `json:"-"`
	Username    string     `json:"username"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Country     *string    `json:"country,omitempty"`
	City        *string    `json:"city,omitempty"`
}

type UserAuthToken struct {
	Token string `json:"token"`
}

type IUserService interface {
	CreateUser(newUser *UserCreateInput) (*User, error)
	CheckPassword(user *User, password string) error
	CreateAuthToken(user *User) (*UserAuthToken, error)
	ReadByEmail(email string) (*User, error)
	GetUser(id string) (*User, error)
}
