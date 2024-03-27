package loadorder

import (
	"context"
	"encoding/json"
	"net/http"
	"skripsi-be/lib/httpx"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
)

func (h eventHandler) HandleLoadOrderEvent(msg *kafka.Message) error {
	var req []params.ServiceInsertOrderToLongTermParam
	err := json.Unmarshal(msg.Value, &req)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.orderService.InsertToLongTerm(context.Background(), req)
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (h httpHandler) LoadOrder(c echo.Context) error {
	var req []params.ServiceInsertOrderToLongTermParam
	err := c.Bind(&req)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	err = h.orderService.InsertToLongTerm(context.Background(), req)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	return httpx.WriteResponse(c, http.StatusCreated, nil)
}
