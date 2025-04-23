package validation_test

import (
	"reflect"
	"testing"

	"beetle/internal/validation"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Age      int    `validate:"gte=0,lte=130"`
	Password string `validate:"required,min=8"`
}

func TestNew(t *testing.T) {
	v := validation.New()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Translator)
}

func TestValidate(t *testing.T) {
	v := validation.New()

	tests := []struct {
		name     string
		input    TestStruct
		hasError bool
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      25,
				Password: "password123",
			},
			hasError: false,
		},
		{
			name: "missing required fields",
			input: TestStruct{
				Name:     "",
				Email:    "",
				Age:      25,
				Password: "",
			},
			hasError: true,
		},
		{
			name: "invalid email",
			input: TestStruct{
				Name:     "John Doe",
				Email:    "invalid-email",
				Age:      25,
				Password: "password123",
			},
			hasError: true,
		},
		{
			name: "invalid age",
			input: TestStruct{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      150,
				Password: "password123",
			},
			hasError: true,
		},
		{
			name: "password too short",
			input: TestStruct{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      25,
				Password: "short",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := v.Validate(tt.input)
			if tt.hasError {
				assert.NotNil(t, errors)
				assert.Greater(t, len(*errors), 0)
			} else {
				assert.Nil(t, errors)
			}
		})
	}
}

func TestTranslateFormErrors(t *testing.T) {
	v := validation.New()
	trans, _ := v.Translator.GetTranslator("en")

	// Create some test errors
	fieldErrors := validation.FormErrors{
		validation.FieldErrors{
			Field: "Name",
			Errors: []validator.FieldError{
				&mockFieldError{field: "Name", tag: "required"},
			},
		},
		validation.FieldErrors{
			Field: "Email",
			Errors: []validator.FieldError{
				&mockFieldError{field: "Email", tag: "email"},
			},
		},
	}

	translated := v.TranslateFormErrors(trans, fieldErrors)
	assert.Equal(t, 2, len(translated))
	assert.Equal(t, "Name", translated[0].Field)
	assert.Equal(t, "Email", translated[1].Field)
	assert.Greater(t, len(translated[0].Errors), 0)
	assert.Greater(t, len(translated[1].Errors), 0)
}

// mockFieldError implements validator.FieldError for testing
type mockFieldError struct {
	field string
	tag   string
}

func (e *mockFieldError) Tag() string {
	return e.tag
}

func (e *mockFieldError) ActualTag() string {
	return e.tag
}

func (e *mockFieldError) Namespace() string {
	return e.field
}

func (e *mockFieldError) StructNamespace() string {
	return e.field
}

func (e *mockFieldError) Field() string {
	return e.field
}

func (e *mockFieldError) StructField() string {
	return e.field
}

func (e *mockFieldError) Value() interface{} {
	return nil
}

func (e *mockFieldError) Param() string {
	return ""
}

func (e *mockFieldError) Kind() reflect.Kind {
	return reflect.String
}

func (e *mockFieldError) Type() reflect.Type {
	return reflect.TypeOf("")
}

func (e *mockFieldError) Translate(ut ut.Translator) string {
	return "translated error"
}

func (e *mockFieldError) Error() string {
	return "error"
}
