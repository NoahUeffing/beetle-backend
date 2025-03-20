package middleware

import "beetle/internal/handler"

type MiddlewareFunc func(*handler.Context) error
