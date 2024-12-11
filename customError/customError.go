package customError

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	Message string `json:"message"`
	Causes  string `json:"causes"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
}

type CustomErrorOption func(customError *CustomError)

func (e CustomError) Error() string {
	return fmt.Sprintf("message: %s - status: %d - causes: %s", e.Message, e.Status, e.Causes)
}

func NewCustomError(opts ...CustomErrorOption) *CustomError {
	err := &CustomError{
		Causes:  "",
		Status:  0,
		Message: "",
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

func WithNotFound(causes, message string) CustomErrorOption {
	return func(e *CustomError) {
		e.Causes = causes
		e.Status = http.StatusNotFound
		e.Message = message
		e.Code = http.StatusText(http.StatusNotFound)
	}
}

func WithBadRequest(causes, message string) CustomErrorOption {
	return func(e *CustomError) {
		e.Causes = causes
		e.Status = http.StatusBadRequest
		e.Message = message
		e.Code = http.StatusText(http.StatusBadRequest)
	}
}

func WithInternalServerError(causes, message string) CustomErrorOption {
	return func(e *CustomError) {
		e.Causes = causes
		e.Status = http.StatusInternalServerError
		e.Message = message
		e.Code = http.StatusText(http.StatusInternalServerError)
	}
}

func WithUnauthorized(causes, message string) CustomErrorOption {
	return func(e *CustomError) {
		e.Causes = causes
		e.Status = http.StatusUnauthorized
		e.Message = message
		e.Code = http.StatusText(http.StatusUnauthorized)
	}
}

func WithCustomError(status int, causes, message string) CustomErrorOption {
	return func(e *CustomError) {
		e.Causes = causes
		e.Status = status
		e.Message = message
		e.Code = http.StatusText(status)
	}
}
