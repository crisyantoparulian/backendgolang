package httphelper_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	httphelper "github.com/SawitProRecruitment/UserService/utils/http_helper"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHttpRespError(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		inputErr     error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "AppError with specific status code",
			inputErr:     apperror.WrapWithCode(errors.New("Bad request"), http.StatusBadRequest),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Bad request"}`,
		},
		{
			name:         "General error returns 500",
			inputErr:     errors.New("Internal Server Error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := httphelper.HttpRespError(c, tt.inputErr) // Explicit package usage

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.JSONEq(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}
