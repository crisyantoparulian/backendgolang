package apperror_test

import (
	"errors"
	"net/http"
	"testing"

	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	errMsg := "something went wrong"
	appErr := apperror.WrapWithCode(errors.New(errMsg), http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, appErr.Code, "HTTP status code should be 400")
	assert.Equal(t, errMsg, appErr.Error(), "Error message should match the original error")
}

func TestWrapWithCode(t *testing.T) {
	errMsg := "database connection failed"
	wrappedErr := apperror.WrapWithCode(errors.New(errMsg), http.StatusInternalServerError)

	assert.NotNil(t, wrappedErr, "Wrapped error should not be nil")
	assert.Equal(t, http.StatusInternalServerError, wrappedErr.Code, "HTTP status should be 500")
	assert.Equal(t, errMsg, wrappedErr.Error(), "Error message should match")
}
