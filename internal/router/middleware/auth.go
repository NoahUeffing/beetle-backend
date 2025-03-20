package middleware

import (
	"beetle/internal/handler"
)

func Auth(hc *handler.Context) error {
	auth, err := hc.AuthService.GetClaims(hc.Context)
	if err != nil {
		return err
	}

	hc.User, err = hc.UserService.ReadByID(auth.UserId)
	return err
}
