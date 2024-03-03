package errors

import (
	"errors"
	"fmt"

	errorsx "github.com/go-errors/errors"
)

var (
	ErrInternalServer                        = errors.New("internal server error")
	ErrBadRequest                            = errors.New("bad request")
	ErrNotFound                              = errors.New("not found")
	ErrUnauthorized                          = errors.New("unauthorized")
	ErrAuthTokenExpired                      = errors.New("token expired")
	ErrIncorrectPassword                     = errors.New("incorrect password")
	ErrDataParamMustNotAFterCurrentTime      = errors.New("data param must not after current time")
	ErrDataParamStartDateMustNotAfterEndDate = errors.New("start date must not after end date")
)

var (
	ErrAccountNotFound         = errors.New("Wrong username or password")
	ErrRecordNotFound          = errors.New("record not found")
	ErrPermissionExpired       = errors.New("permission expired")
	ErrUsernameDuplicated      = errors.New("username is already exist")
	ErrUsernameNotExist        = errors.New("username not exists")
	ErrBeginTransaction        = errors.New("failed begin transaction")
	ErrCommitTransaction       = errors.New("failed commit transaction")
	ErrCustomerCashierNotMatch = errors.New("customer and cashier not match")
	ErrStoreNotFound           = errors.New("store not found")
	ErrCustomerNotFound        = errors.New("customer not found")
	ErrCashierNotFound         = errors.New("cashier not found")
	ErrItemNotFound            = errors.New("item not found")
)

var (
	ErrValidation = errorsx.New("validation error")
)

var (
	ErrJWTMissingOrInvalid = errors.New("jwt token missing or invalid")
	ErrFailedCastJWTClaims = errors.New("failed to cast jwt claims")
)

func Join(errs ...error) error {
	return errorsx.Join(errs...)
}

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
	if err != nil {
		var errorsErr *errorsx.Error
		if errorsx.As(err, &errorsErr) {
			fmt.Println(errorsErr.ErrorStack())
		}
	}
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
