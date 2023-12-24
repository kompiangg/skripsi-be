package httpx

import (
	"errors"
	"log"
	"net/http"
	"strings"

	x "skripsi-be/pkg/errors"
	httppkg "skripsi-be/pkg/http"

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
	e := httppkg.GetResponseErr(errParam)

	if x.Is(errParam, x.ErrValidation) {
		e.Message = x.ErrBadRequest.Error()
		e.HTTPErrorCode = echo.ErrBadRequest.Code

		// To getting the Unwrap method from private object joinError in "errors" package
		var joinErr interface{ Unwrap() []error }
		if errors.As(errParam, &joinErr) {
			errs := joinErr.Unwrap()[1].Error()
			detail = strings.Split(errs, "\n --- ")[1:]
		}
	} else if e.HTTPErrorCode == http.StatusInternalServerError {
		x.ErrorStack(errParam)
		detail = nil
	} else {
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
