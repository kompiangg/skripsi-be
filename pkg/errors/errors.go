package errors

import (
	"errors"

	errorsx "github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

var (
	ErrInternalServer    = errors.New("internal server error")
	ErrBadRequest        = errors.New("bad request")
	ErrNotFound          = errors.New("not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrAuthTokenExpired  = errors.New("token expired")
	ErrIncorrectPassword = errors.New("incorrect password")
)

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrRecordNotFound     = errors.New("record not found")
	ErrPermissionExpired  = errors.New("permission expired")
	ErrUsernameDuplicated = errors.New("username is already exist")
	ErrUsernameNotExist   = errors.New("username not exists")
	ErrBeginTransaction   = errors.New("failed begin transaction")
	ErrCommitTransaction  = errors.New("failed commit transaction")
)

var (
	ErrValidation = errorsx.New("validation error")
)

var (
	ErrJWTMissingOrInvalid = errors.New("jwt token missing or invalid")
	ErrFailedCastJWTClaims = errors.New("failed to cast jwt claims")
)

func Wrap(cause error) error {
	if cause == nil {
		return nil
	}

	return errorsx.Wrap(cause, 0)
}

func Unwrap(err error) error {
	if err != nil {
		var errorsErr *errorsx.Error
		if errorsx.As(err, &errorsErr) {
			return errorsErr.Unwrap()
		}
	}

	return err
}

func ErrorStack(err error) {
	log.Err(err).Msgf("%+v", err.(*errorsx.Error).ErrorStack())
}

func New(e interface{}) *errorsx.Error {
	return errorsx.New(e)
}

func As(err error, target any) bool {
	return errorsx.As(err, target)
}

func Is(err error, target error) bool {
	return errorsx.Is(err, target)
}
