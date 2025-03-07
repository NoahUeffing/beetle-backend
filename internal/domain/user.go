package domain

import "time"

type UserAuthInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserCreateInput struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type User struct {
	ID          string     `json:"id"`
	Email       string     `json:"email" validate:"email"`
	Password    string     `json:"-"`
	Username    string     `json:"username"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Country     *string    `json:"country,omitempty"`
	City        *string    `json:"city,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type IUserService interface {
	CreateUser(newUser *UserCreateInput) (*User, error)
	CheckPassword(user *User, password string) error
}
