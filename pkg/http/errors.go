package http

import (
	"net/http"

	x "skripsi-be/pkg/errors"
)

type ErrorSchema struct {
	HTTPErrorCode int
	Message       string
}

var errMap map[error]ErrorSchema = map[error]ErrorSchema{
	x.ErrInternalServer:                        {HTTPErrorCode: http.StatusInternalServerError, Message: x.ErrInternalServer.Error()},
	x.ErrBadRequest:                            {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrBadRequest.Error()},
	x.ErrValidation:                            {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrBadRequest.Error()},
	x.ErrRecordNotFound:                        {HTTPErrorCode: http.StatusNotFound, Message: x.ErrRecordNotFound.Error()},
	x.ErrAccountNotFound:                       {HTTPErrorCode: http.StatusNotFound, Message: x.ErrAccountNotFound.Error()},
	x.ErrNotFound:                              {HTTPErrorCode: http.StatusNotFound, Message: x.ErrNotFound.Error()},
	x.ErrUnauthorized:                          {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrUnauthorized.Error()},
	x.ErrPermissionExpired:                     {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrPermissionExpired.Error()},
	x.ErrAuthTokenExpired:                      {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrAuthTokenExpired.Error()},
	x.ErrUsernameDuplicated:                    {HTTPErrorCode: http.StatusConflict, Message: x.ErrUsernameDuplicated.Error()},
	x.ErrIncorrectPassword:                     {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrIncorrectPassword.Error()},
	x.ErrUsernameNotExist:                      {HTTPErrorCode: http.StatusNotFound, Message: x.ErrUsernameNotExist.Error()},
	x.ErrDataParamMustNotAFterCurrentTime:      {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrDataParamMustNotAFterCurrentTime.Error()},
	x.ErrDataParamStartDateMustNotAfterEndDate: {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrDataParamStartDateMustNotAfterEndDate.Error()},
	x.ErrCustomerCashierNotMatch:               {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrCustomerCashierNotMatch.Error()},
	x.ErrStoreNotFound:                         {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrStoreNotFound.Error()},
	x.ErrCustomerNotFound:                      {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrCustomerNotFound.Error()},
	x.ErrCashierNotFound:                       {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrCashierNotFound.Error()},
	x.ErrItemNotFound:                          {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrItemNotFound.Error()},
}

func GetResponseErr(param error) ErrorSchema {
	param = x.Unwrap(param)

	res, exists := errMap[param]
	if !exists {
		return errMap[x.ErrInternalServer]
	}

	return res
}
