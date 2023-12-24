package errors

import (
	"errors"
	"os"

	errorsx "github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
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
	ErrAccountNotFound   = errors.New("account not found")
	ErrDoctorNotFound    = errors.New("doctor not found")
	ErrRecordNotFound    = errors.New("record not found")
	ErrPermissionExpired = errors.New("permission expired")
	ErrEmailDuplicated   = errors.New("email is already exist")
	ErrUsernameNotExist  = errors.New("username not exists")
	ErrBeginTransaction  = errors.New("failed begin transaction")
	ErrCommitTransaction = errors.New("failed commit transaction")
)

var (
	ErrValidation = errorsx.New("validation error")
)

var (
	ErrJWTMissingOrInvalid = errors.New("jwt token missing or invalid")
	ErrFailedCastJWTClaims = errors.New("failed to cast jwt claims")
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		DisableQuote:  true,
		PadLevelText:  false,
	})
}

func Wrap(cause error, msg string) error {
	if cause == nil {
		return nil
	}

	return errorsx.WrapPrefix(cause, msg, 0)
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
	log.Warningln(errorsx.New(err).ErrorStack())
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
