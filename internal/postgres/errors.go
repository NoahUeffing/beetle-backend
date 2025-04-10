package postgres

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrEntityVersionConflict  = errors.New("entity submitted is out of date")
	ErrEntityNotFound         = errors.New("entity was not found")
	ErrEntityNonUnique        = errors.New("the entity you are trying to create already exists")
	ErrInvalidCredentials     = errors.New("email or password invalid")
	ErrEmailAlreadyAssociated = errors.New("email is already associated with user")
	ErrUsernameTaken          = errors.New("username is not available")
	ErrProfaneString          = errors.New("non-profane string could not be created")
	ErrReadingEntity          = errors.New("error reading entity")
	ErrInvalidAction          = errors.New("invalid action")
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
