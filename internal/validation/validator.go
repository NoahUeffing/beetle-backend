package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Need to alias to allow for `Validate` method naming
// type gpValidator = validator.Validate

type Validator struct {
	gpValidator validator.Validate
	Translator  *ut.UniversalTranslator
}

type FieldErrors struct {
	Field  string
	Errors []validator.FieldError
}

type TranslatedFieldErrors struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

type FormErrors = []FieldErrors
type TranslatedFormErrors = []TranslatedFieldErrors

func (v *Validator) initTranslator() {
	if v.Translator != nil {
		return
	}

	en := en.New()
	// Note: Other languages go here
	v.Translator = ut.New(en, en)
	trans, _ := v.Translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(&v.gpValidator, trans)
}

func New() *Validator {
	v := &Validator{
		gpValidator: *validator.New(),
	}
	v.initTranslator()

	return v
}

func convertErrors(err error) *[]FieldErrors {
	if err == nil {
		return nil
	}

	errorList, _ := err.(validator.ValidationErrors)
	if len(errorList) == 0 {
		return nil
	}

	errMap := map[string][]validator.FieldError{}
	for _, e := range errorList {
		field := e.Field()
		errMap[field] = append(errMap[field], e)
	}

	result := make([]FieldErrors, 0)
	for k, v := range errMap {
		result = append(result, FieldErrors{
			Field:  k,
			Errors: v,
		})
	}

	return &result
}

func (v *Validator) Validate(i interface{}) *FormErrors {
	err := v.gpValidator.Struct(i)
	if err == nil {
		return nil
	}
	return convertErrors(err.(validator.ValidationErrors))
}

func TranslateFieldErrors(t ut.Translator, fe FieldErrors) TranslatedFieldErrors {
	tfe := TranslatedFieldErrors{
		Field:  fe.Field,
		Errors: []string{},
	}

	for _, e := range fe.Errors {
		tfe.Errors = append(tfe.Errors, e.Translate(t))
	}

	return tfe
}

func (v *Validator) TranslateFormErrors(t ut.Translator, list FormErrors) TranslatedFormErrors {
	result := make(TranslatedFormErrors, len(list))
	for i, fe := range list {
		result[i] = TranslateFieldErrors(t, fe)
	}

	return result
}
