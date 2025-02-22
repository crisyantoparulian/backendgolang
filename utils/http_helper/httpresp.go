package httphelper

import (
	"net/http"

	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func HttpRespError(c echo.Context, srcErr error) (err error) {
	// Check if it's an APIError
	if apiErr, ok := srcErr.(*apperror.AppError); ok {
		err = c.JSON(apiErr.Code, ErrorResponse{
			Message: srcErr.Error(),
		})
	} else {
		err = c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: srcErr.Error(),
		})
	}
	return
}
