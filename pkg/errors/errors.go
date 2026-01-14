package errors

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details,omitempty"`
	HTTPStatus int               `json:"-"`
	GRPCCode   codes.Code        `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) GRPCStatus() *status.Status {
	return status.New(e.GRPCCode, e.Message)
}

var (
	ErrInvalidArgument = &AppError{
		Code:       "INVALID_ARGUMENT",
		Message:    "Invalid argument",
		HTTPStatus: http.StatusBadRequest,
		GRPCCode:   codes.InvalidArgument,
	}

	ErrNotFound = &AppError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
		GRPCCode:   codes.NotFound,
	}

	ErrAlreadyExists = &AppError{
		Code:       "ALREADY_EXISTS",
		Message:    "Resource already exists",
		HTTPStatus: http.StatusConflict,
		GRPCCode:   codes.AlreadyExists,
	}

	ErrUnauthenticated = &AppError{
		Code:       "UNAUTHENTICATED",
		Message:    "Authentication required",
		HTTPStatus: http.StatusUnauthorized,
		GRPCCode:   codes.Unauthenticated,
	}

	ErrPermissionDenied = &AppError{
		Code:       "PERMISSION_DENIED",
		Message:    "Permission denied",
		HTTPStatus: http.StatusForbidden,
		GRPCCode:   codes.PermissionDenied,
	}

	ErrResourceExhausted = &AppError{
		Code:       "RESOURCE_EXHAUSTED",
		Message:    "Resource exhausted",
		HTTPStatus: http.StatusTooManyRequests,
		GRPCCode:   codes.ResourceExhausted,
	}

	ErrInternal = &AppError{
		Code:       "INTERNAL",
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
		GRPCCode:   codes.Internal,
	}

	ErrUnavailable = &AppError{
		Code:       "UNAVAILABLE",
		Message:    "Service unavailable",
		HTTPStatus: http.StatusServiceUnavailable,
		GRPCCode:   codes.Unavailable,
	}
)

// New creates a new AppError with custom message
func New(base *AppError, message string) *AppError {
	return &AppError{
		Code:       base.Code,
		Message:    message,
		HTTPStatus: base.HTTPStatus,
		GRPCCode:   base.GRPCCode,
	}
}

// WithDetails adds details to an error
func WithDetails(base *AppError, details map[string]string) *AppError {
	return &AppError{
		Code:       base.Code,
		Message:    base.Message,
		Details:    details,
		HTTPStatus: base.HTTPStatus,
		GRPCCode:   base.GRPCCode,
	}
}

// FromGRPC converts a gRPC error to AppError
func FromGRPC(err error) *AppError {
	st, ok := status.FromError(err)
	if !ok {
		return ErrInternal
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return New(ErrInvalidArgument, st.Message())
	case codes.NotFound:
		return New(ErrNotFound, st.Message())
	case codes.AlreadyExists:
		return New(ErrAlreadyExists, st.Message())
	case codes.Unauthenticated:
		return New(ErrUnauthenticated, st.Message())
	case codes.PermissionDenied:
		return New(ErrPermissionDenied, st.Message())
	case codes.ResourceExhausted:
		return New(ErrResourceExhausted, st.Message())
	case codes.Unavailable:
		return New(ErrUnavailable, st.Message())
	default:
		return New(ErrInternal, st.Message())
	}
}
