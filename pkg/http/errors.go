package http

import (
	"net/http"

	x "skripsi-be/pkg/errors"
)

type errorSchema struct {
	HTTPErrorCode int
	Message       string
}

var errMap map[error]errorSchema = map[error]errorSchema{
	x.ErrInternalServer:    {HTTPErrorCode: http.StatusInternalServerError, Message: x.ErrInternalServer.Error()},
	x.ErrBadRequest:        {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrBadRequest.Error()},
	x.ErrValidation:        {HTTPErrorCode: http.StatusBadRequest, Message: x.ErrBadRequest.Error()},
	x.ErrRecordNotFound:    {HTTPErrorCode: http.StatusNotFound, Message: x.ErrRecordNotFound.Error()},
	x.ErrAccountNotFound:   {HTTPErrorCode: http.StatusNotFound, Message: x.ErrAccountNotFound.Error()},
	x.ErrDoctorNotFound:    {HTTPErrorCode: http.StatusNotFound, Message: x.ErrDoctorNotFound.Error()},
	x.ErrNotFound:          {HTTPErrorCode: http.StatusNotFound, Message: x.ErrNotFound.Error()},
	x.ErrUnauthorized:      {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrUnauthorized.Error()},
	x.ErrPermissionExpired: {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrPermissionExpired.Error()},
	x.ErrAuthTokenExpired:  {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrAuthTokenExpired.Error()},
	x.ErrEmailDuplicated:   {HTTPErrorCode: http.StatusConflict, Message: x.ErrEmailDuplicated.Error()},
	x.ErrIncorrectPassword: {HTTPErrorCode: http.StatusUnauthorized, Message: x.ErrIncorrectPassword.Error()},
	x.ErrUsernameNotExist:  {HTTPErrorCode: http.StatusNotFound, Message: x.ErrUsernameNotExist.Error()},
}

func GetResponseErr(param error) errorSchema {
	param = x.Unwrap(param)

	res, exists := errMap[param]
	if !exists {
		return errMap[x.ErrInternalServer]
	}

	return res
}
