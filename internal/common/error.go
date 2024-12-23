package common

import (
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("not found")
var ErrEmailAlreadyExists = errors.New("email already exists")

// ConvertErrorToHTTPStatus converts error to HTTP status code and message.
// This is the simplified version of the error handling in the http_handler package.
// Other options are using custom error type, and using middleware to handle error.
func ConvertErrorToHTTPStatus(err error) (int, string) {
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, "not found"
	}
	if errors.Is(err, ErrEmailAlreadyExists) {
		return http.StatusConflict, "email already exists"
	}
	return http.StatusInternalServerError, "internal server error"
}
