package apperror

type AppError struct {
	Code int   // HTTP status code
	err  error // Error message
}

// Implement the `error` interface
func (e *AppError) Error() string {
	return e.err.Error()
}

// Helper functions to wrap http code
func WrapWithCode(err error, code int) *AppError {
	return &AppError{Code: code, err: err}
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}
