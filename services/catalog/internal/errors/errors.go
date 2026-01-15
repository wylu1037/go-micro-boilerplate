package errors

import (
	stderrors "errors"

	microerrors "go-micro.dev/v4/errors"
)

const serviceName = "catalog"

var (
	ErrShowNotFound      = stderrors.New("show not found")
	ErrVenueNotFound     = stderrors.New("venue not found")
	ErrSessionNotFound   = stderrors.New("session not found")
	ErrSeatAreaNotFound  = stderrors.New("seat area not found")
	ErrInsufficientSeats = stderrors.New("insufficient seats available")
	ErrInvalidSeatArea   = stderrors.New("seat area does not belong to session")
	ErrInvalidPrice      = stderrors.New("invalid price format")
)

func ToMicroError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case stderrors.Is(err, ErrShowNotFound):
		return microerrors.NotFound(serviceName, "show not found")
	case stderrors.Is(err, ErrVenueNotFound):
		return microerrors.NotFound(serviceName, "venue not found")
	case stderrors.Is(err, ErrSessionNotFound):
		return microerrors.NotFound(serviceName, "session not found")
	case stderrors.Is(err, ErrSeatAreaNotFound):
		return microerrors.NotFound(serviceName, "seat area not found")
	case stderrors.Is(err, ErrInsufficientSeats):
		return microerrors.BadRequest(serviceName, "insufficient seats available")
	case stderrors.Is(err, ErrInvalidSeatArea):
		return microerrors.BadRequest(serviceName, "seat area does not belong to session")
	case stderrors.Is(err, ErrInvalidPrice):
		return microerrors.BadRequest(serviceName, "invalid price format")
	default:
		return microerrors.InternalServerError(serviceName, "internal server error")
	}
}
