package validator_test

import (
	"testing"

	validatorutils "github.com/SawitProRecruitment/UserService/utils/validator" // Import package explicitly
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Age   int    `validate:"gt=0"`
	Email string `validate:"email"`
}

func TestFormatValidationErrors(t *testing.T) {
	validate := validator.New()
	testData := TestStruct{
		Name:  "",       // Required field (should trigger an error)
		Age:   -5,       // Should be greater than 0 (should trigger an error)
		Email: "wrong@", // Invalid email format (should trigger an error)
	}

	err := validate.Struct(testData)
	if err == nil {
		t.Fatalf("Expected validation errors but got none")
	}

	// Convert validation errors
	validationErrors := err.(validator.ValidationErrors)
	formattedErrors := validatorutils.FormatValidationErrors(validationErrors)

	// Expected error messages
	expectedErrors := map[string]string{
		"Name":  "This field is required",
		"Age":   "Must be greater than 0",
		"Email": "Invalid value",
	}

	assert.Equal(t, expectedErrors, formattedErrors, "Validation error messages should match expected values")
}
