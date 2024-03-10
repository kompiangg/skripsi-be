package service

import (
	"skripsi-be/config"
	"skripsi-be/repository"
	"skripsi-be/service/auth"
	"skripsi-be/service/cashier"
	"skripsi-be/service/customer"
	"skripsi-be/service/item"
	"skripsi-be/service/order"
	"skripsi-be/service/payment_types"
)

type Service struct {
	Order        order.Service
	Auth         auth.Service
	Item         item.Service
	PaymentTypes payment_types.Service
	Customer     customer.Service
	Cashier      cashier.Service
}

func New(
	repository repository.Repository,
	config config.Config,
) (Service, error) {
	order := order.New(
		order.Config{
			IsUsingSharding: config.ShardingDatabase.IsUsingSharding,
			Shards:          config.ShardingDatabase.Shards,
			Date:            config.Date,
		},
		repository.Order,
		repository.Publisher,
		repository.Currency,
		repository.Cashier,
		repository.Customer,
		repository.Store,
		repository.Item,
		repository.Scheduler,

		repository.ShardDBTx,
		repository.LongTermDBTx,
		repository.GeneralDBTx,
		repository.GetShardIndexByDateTime,
		repository.GetShardWhereQuery,
	)

	auth := auth.New(
		auth.Config{
			JWT: config.JWT,
		},
		repository.Account,
		repository.Admin,
		repository.Cashier,
	)

	item := item.New(
		item.Config{},
		repository.Item,
	)

	paymentTypes := payment_types.New(
		payment_types.Config{},
		repository.PaymentTypes,
	)

	customer := customer.New(
		customer.Config{},
		repository.Customer,
	)

	cashier := cashier.New(
		cashier.Config{},
		repository.Cashier,
	)

	return Service{
		Order:        order,
		Auth:         auth,
		Item:         item,
		PaymentTypes: paymentTypes,
		Customer:     customer,
		Cashier:      cashier,
	}, nil
}
