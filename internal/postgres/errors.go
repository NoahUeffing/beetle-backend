package postgres

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrEntityNotFound     = errors.New("entity was not found")
	ErrEntityNonUnique    = errors.New("the entity you are trying to create already exists")
	ErrInvalidCredentials = errors.New("email or password invalid")
)

func ConvertErrorIfNeeded(err error) error {
	if err == sql.ErrNoRows {
		return ErrEntityNotFound
	}

	var e *pq.Error
	if errors.As(err, &e) {
		switch e.Code.Name() { // See https://www.postgresql.org/docs/9.3/errcodes-appendix.html
		case "unique_violation":
			return ErrEntityNonUnique
		default:
			return err
		}
	}

	return err
}
