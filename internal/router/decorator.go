package router

import (
	"beetle/internal/handler"
	"beetle/internal/router/middleware"

	"github.com/labstack/echo/v4"
)

type handlerFunc func(handler.Context) error

func WithContext(hf handlerFunc, mfs ...middleware.MiddlewareFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sc := c.(*handler.Context)
		var err error
		for _, mf := range mfs {
			err = mf(sc)
			if err != nil {
				return err
			}
		}
		return hf(*sc)
	}
}
