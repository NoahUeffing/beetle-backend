package handler

// A simple wrapper for validation which validates and translates returned errors
func (c *Context) validate(i interface{}) *FormValidationError {
	formErrs := c.ValidationService.Validate(i)

	if formErrs == nil || len(*formErrs) == 0 {
		return nil
	}

	// Note: translator should eventually be chosen by AcceptLanguage Header from context eventually
	t := c.ValidationService.Translator.GetFallback()
	return &FormValidationError{
		Fields: c.ValidationService.TranslateFormErrors(t, *formErrs),
	}
}
