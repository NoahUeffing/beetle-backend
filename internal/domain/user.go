package domain

type User struct {
}

type UserAuthInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type IUserService interface {
}
