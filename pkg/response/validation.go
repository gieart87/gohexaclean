package response

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ParseValidationErrors converts ozzo-validation errors to array format
func ParseValidationErrors(err error) map[string][]string {
	errors := make(map[string][]string)

	if validationErrs, ok := err.(validation.Errors); ok {
		for field, fieldErr := range validationErrs {
			// Each field gets an array of error messages
			errors[field] = []string{fieldErr.Error()}
		}
	} else {
		// If not a validation error, return as general error
		errors["error"] = []string{err.Error()}
	}

	return errors
}
