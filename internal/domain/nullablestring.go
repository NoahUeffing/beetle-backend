package domain

import (
	"database/sql/driver"
	"errors"
)

type NullableString string

// Implementing https://pkg.go.dev/database/sql#Scanner
func (s *NullableString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return errors.New("attempted to scan a non-string column into a NullableString")
	}
	*s = NullableString(strVal)
	return nil
}

// Implementing https://pkg.go.dev/database/sql/driver#Valuer
func (s NullableString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}
