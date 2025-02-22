package validator

import "github.com/go-playground/validator/v10"

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errorsMap := make(map[string]string)
	for _, err := range errs {
		errorsMap[err.Field()] = getValidationMessage(err)
	}
	return errorsMap
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "gt":
		return "Must be greater than 0"
	default:
		return "Invalid value"
	}
}
