package nuxeogoclient

import "fmt"

type RequestError struct {
	StatusCode int
}

func newRequestError(statusCode int) *RequestError {
	e := new(RequestError)
	e.StatusCode = statusCode
	return e
}

func (e RequestError) Error() string {
	return fmt.Sprintf("Request error. Status code: %d", e.StatusCode)
}

func (e RequestError) IsUserError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

func (e RequestError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}
