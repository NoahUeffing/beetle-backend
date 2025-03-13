package handler

import (
	"beetle/internal/validation"
)

func translateErrors(v *validation.Validator, i interface{}, errors *validation.FormErrors) map[string]interface{} {
	trans, _ := v.Translator.GetTranslator("en")

	translatedErrors := v.TranslateFormErrors(trans, *errors)

	return map[string]interface{}{
		"error":  "Validation failed",
		"errors": translatedErrors,
	}
}
