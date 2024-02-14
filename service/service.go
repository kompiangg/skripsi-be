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
		repository.Item,

		repository.ShardDBTx,
		repository.LongTermDBTx,
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

	return Service{
		Order: order,
		Auth:  auth,
	}, nil
}
