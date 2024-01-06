package service

import (
	"skripsi-be/config"
	"skripsi-be/repository"
	"skripsi-be/service/order"
)

type Service struct {
	Order order.Order
}

func New(
	repository repository.Repository,
	config config.Config,
) (Service, error) {
	order := order.New(
		order.Config{
			Shards: config.ShardingDatabase.Shards,
		},
		repository.Sharding,
		repository.Order,
		repository.BeginShardDBTx,
	)

	return Service{
		Order: order,
	}, nil
}
