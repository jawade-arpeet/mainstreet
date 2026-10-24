package errs

import (
	"fmt"
	"net/http"
)

type Code string

const (
	CodeInternal Code = "INTERNAL_SERVER_ERROR"
)

type HTTPError struct {
	Code    Code   `json:"code"`
	Status  int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func newHTTPError(code Code, status int, msg string, err error) *HTTPError {
	return &HTTPError{Code: code, Status: status, Message: msg, Err: err}
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func BadRequest(code Code, msg string, err error) *HTTPError {
	return newHTTPError(code, http.StatusBadRequest, msg, err)
}

func Unauthorized(code Code, msg string, err error) *HTTPError {
	return newHTTPError(code, http.StatusUnauthorized, msg, err)
}

func NotFound(code Code, msg string, err error) *HTTPError {
	return newHTTPError(code, http.StatusNotFound, msg, err)
}

func Conflict(code Code, msg string, err error) *HTTPError {
	return newHTTPError(code, http.StatusConflict, msg, err)
}

func Internal(err error) *HTTPError {
	return newHTTPError(
		CodeInternal,
		http.StatusInternalServerError,
		"something went wrong",
		err,
	)
}
