package errors

import (
	stderrors "errors"

	microerrors "go-micro.dev/v4/errors"
)

const serviceName = "identity"

var (
	ErrUserNotFound       = stderrors.New("user not found")
	ErrUserAlreadyExists  = stderrors.New("user already exists")
	ErrInvalidCredentials = stderrors.New("invalid credentials")
)

var (
	ErrTokenNotFound = stderrors.New("token not found")
	ErrTokenExpired  = stderrors.New("token expired")
	ErrTokenUsed     = stderrors.New("token already used")
)

func ToMicroError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case stderrors.Is(err, ErrUserNotFound):
		return microerrors.NotFound(serviceName, "user not found")
	case stderrors.Is(err, ErrUserAlreadyExists):
		return microerrors.Conflict(serviceName, "user already exists")
	case stderrors.Is(err, ErrInvalidCredentials):
		return microerrors.Unauthorized(serviceName, "invalid credentials")
	case stderrors.Is(err, ErrTokenNotFound):
		return microerrors.Unauthorized(serviceName, "invalid token")
	case stderrors.Is(err, ErrTokenExpired):
		return microerrors.Unauthorized(serviceName, "token expired")
	case stderrors.Is(err, ErrTokenUsed):
		return microerrors.BadRequest(serviceName, "token already used")
	default:
		return microerrors.InternalServerError(serviceName, "internal server error")
	}
}
