package errres

import (
	"net/http"
)

var (
	ErrDuplicateResource     = NewWithCode(http.StatusConflict, "duplicated resource")
	ErrInvalidRequest        = NewWithCode(http.StatusBadRequest, "invalid request ")
	ErrResourceNotFound      = NewWithCode(http.StatusNotFound, "resource not found")
	ErrAuthorizationRequired = NewWithCode(http.StatusUnauthorized, "authorization required")
	ErrEndpointNotFound      = NewWithCode(http.StatusNotFound, "endpoint not found")
	ErrFailedRequest         = NewWithCode(http.StatusExpectationFailed, "failed request ")
)

type HttpError struct {
	Message    string `json:"message"`
	Code       int    `json:"code"`
	InnerError error
}

func (he *HttpError) Error() string {
	return he.Message
}

var _ error = &HttpError{}

func New(message string) *HttpError {
	return NewWithCode(http.StatusInternalServerError, message)
}

func NewWithCode(code int, message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    code,
	}
}

func Wrap(message string, err error) *HttpError {
	return WrapWithCode(http.StatusInternalServerError, message, err)
}

func WrapWithCode(code int, message string, err error) *HttpError {
	erro := NewWithCode(code, message)
	erro.InnerError = err
	return erro
}

func BadRequest(message string, err error) *HttpError {
	erro := NewWithCode(http.StatusBadRequest, message)
	erro.InnerError = err
	return erro
}
