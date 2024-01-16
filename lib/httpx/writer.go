package httpx

import (
	"log"
	"net/http"

	x "skripsi-be/pkg/errors"
	httppkg "skripsi-be/pkg/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

func WriteResponse(c echo.Context, code int, data interface{}) error {
	if data == nil {
		data = http.StatusText(code)
	}

	err := c.JSON(code, HTTPBaseResponse{
		Error: nil,
		Data:  data,
	})
	if err != nil {
		log.Println("[WriteResponse] FATAL ERROR on send response to client:", err)
		return err
	}

	return nil
}

func WriteErrorResponse(c echo.Context, errParam error, detail interface{}) error {
	var e httppkg.ErrorSchema
	var echoHTTPError *echo.HTTPError
	var validationError validation.Errors

	if x.As(errParam, &validationError) {
		e.Message = x.ErrBadRequest.Error()
		e.HTTPErrorCode = echo.ErrBadRequest.Code
		detail = validationError.Error()
	} else if x.As(errParam, &echoHTTPError) {
		if echoHTTPError.Code == http.StatusInternalServerError {
			x.ErrorStack(errParam)
			detail = nil
		} else if echoHTTPError.Code == http.StatusBadRequest {
			e.Message = x.ErrBadRequest.Error()
			e.HTTPErrorCode = echo.ErrBadRequest.Code
			detail = echoHTTPError.Error()
		}
	} else {
		e = httppkg.GetResponseErr(errParam)
		if e.HTTPErrorCode == http.StatusInternalServerError {
			x.ErrorStack(errParam)
		}

		detail = nil
	}

	err := c.JSON(e.HTTPErrorCode, HTTPBaseResponse{
		Error: &HTTPErrorBaseResponse{
			Message: e.Message,
			Detail:  detail,
		},
		Data: nil,
	})

	if err != nil {
		log.Println("[WriteErrorResponse] FATAL ERROR on send response to client:", err)
		return err
	}

	return nil
}
