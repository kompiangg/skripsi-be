package service

import (
	"skripsi-be/config"
	"skripsi-be/repository"
	"skripsi-be/service/auth"
	"skripsi-be/service/order"
)

type Service struct {
	Order order.Service
	Auth  auth.Service
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
		repository.Sharding,
		repository.Order,
		repository.Publisher,
		repository.Currency,
		repository.Cashier,
		repository.Customer,
		repository.Store,

		repository.ShardDBTx,
		repository.LongTermDBTx,
		repository.GetShardIndexByDateTime,
		repository.GetShardWhereQuery,
	)

	auth := auth.New(
		auth.Config{
			Admin:   config.JWT.Admin,
			Cashier: config.JWT.Cashier,
		},
		repository.Account,
		repository.Admin,
	)

	return Service{
		Order: order,
		Auth:  auth,
	}, nil
}
