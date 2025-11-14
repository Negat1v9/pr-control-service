package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Negat1v9/pr-review-service/pkg/logger"
)

const (
	ErrBadRequest      = "BAD_REQUEST"
	ErrNotFound        = "NOT_FOUND"
	ErrPrExists        = "PR_EXISTS"
	ErrTeamExists      = "TEAM_EXISTS"
	ErrRequestTimeout  = "REQUEST_TIMEOUT"
	ErrInternal        = "INTERNAL_SERVER_ERROR"
	ErrUserNotReviewer = "NOT_ASSIGNED"
	ErrNoCantidate     = "NO_CANDIDATE"
	ErrPrAlredyMerged  = "PR_MERGED"
)

type Error struct {
	StatusCode int    `json:"-"`       // http status code
	Code       string `json:"code"`    // error code
	Message    string `json:"message"` // error message
	Causes     any    `json:"-"`       // error causes for internal use
}

// Errors - implementation of the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("status: %d, message: %s, causes: %v", e.StatusCode, e.Message, e.Causes)
}

// New creates a new HTTP error
func NewError(statusCode int, code string, message string, causes any) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Causes:     causes,
	}
}

func NewBadRequestError(message string, causes any) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Code:       ErrBadRequest,
		Message:    message,
		Causes:     causes,
	}
}

func NewNotFoundError(message string, causes any) *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Code:       ErrNotFound,
		Message:    message,
		Causes:     causes,
	}
}

// ParseError - parses an error into an HTTP error
func parseError(err error) *Error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewError(http.StatusNotFound, "resource not found", ErrNotFound, nil)
	case errors.Is(err, context.DeadlineExceeded):
		return NewError(http.StatusRequestTimeout, "request timeout", ErrRequestTimeout, nil)
	case strings.Contains(err.Error(), "Unmarshal"):
		return NewError(http.StatusBadRequest, "bad request", ErrBadRequest, nil)
	default:
		if restErr, ok := err.(*Error); ok {
			return restErr
		}
		return NewError(http.StatusInternalServerError, "internal server error", ErrInternal, nil)
	}
}

func LogResponseErr(r *http.Request, log *logger.Logger, err error) {
	if err != nil {
		log.Errorf("Path: %s, Error: %s", r.RequestURI, err.Error())
	}
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, dataName string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if dataName == "" {
		json.NewEncoder(w).Encode(data)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{dataName: data})
}

func WriteErrResponse(w http.ResponseWriter, err error) {
	httpErr := parseError(err)

	WriteJsonResponse(w, httpErr.StatusCode, "error", httpErr)
}
