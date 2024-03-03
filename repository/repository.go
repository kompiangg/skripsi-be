package repository

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/repository/account"
	"skripsi-be/repository/admin"
	"skripsi-be/repository/cashier"
	"skripsi-be/repository/currency"
	"skripsi-be/repository/customer"
	"skripsi-be/repository/item"
	"skripsi-be/repository/order"
	"skripsi-be/repository/payment_types"
	"skripsi-be/repository/publisher"
	"skripsi-be/repository/shard"
	"skripsi-be/repository/store"
	"skripsi-be/type/result"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Sharding     shard.Repository
	Order        order.Repository
	Account      account.Repository
	Admin        admin.Repository
	Publisher    publisher.Repository
	Currency     currency.Repository
	Cashier      cashier.Repository
	Customer     customer.Repository
	Store        store.Repository
	Item         item.Repository
	PaymentTypes payment_types.Repository

	LongTermDBTx            func(ctx context.Context) (*sqlx.Tx, error)
	ShardDBTx               func(ctx context.Context, dbIndex int) (*sqlx.Tx, error)
	GetShardIndexByDateTime func(date time.Time) (int, error)
	GetShardWhereQuery      func(startDate time.Time, endDate time.Time) ([]result.ShardTimeSeriesWhereQuery, error)
}

func New(
	config config.Config,
	longTermDatabase *sqlx.DB,
	generalDatabase *sqlx.DB,
	shardingDatabase []*sqlx.DB,
	redis *redis.Client,
	kafkaPublisher *kafka.Producer,
) (Repository, error) {
	sharding := shard.New(
		shard.Config{
			Shards: config.ShardingDatabase,
		},
		redis,
	)

	order := order.New(
		order.Config{
			ShardingDatabase: config.ShardingDatabase,
		},
		shardingDatabase,
		longTermDatabase,
	)

	account := account.New(
		account.Config{},
		generalDatabase,
	)

	admin := admin.New(
		admin.Config{},
		generalDatabase,
	)

	publisher := publisher.New(
		publisher.Config{
			LoadOrderTopic:      config.Kafka.Topic.LoadOrder,
			TransformOrderTopic: config.Kafka.Topic.TransformOrder,
		},
		kafkaPublisher,
	)

	currency := currency.New(
		currency.Config{},
		redis,
	)

	cashier := cashier.New(
		cashier.Config{},
		generalDatabase,
	)

	customer := customer.New(
		customer.Config{},
		generalDatabase,
	)

	store := store.New(
		store.Config{},
		generalDatabase,
	)

	item := item.New(
		item.Config{},
		generalDatabase,
	)

	paymentTypes := payment_types.New(
		payment_types.Config{},
		generalDatabase,
	)

	return Repository{
		Sharding:     sharding,
		Order:        order,
		Account:      account,
		Admin:        admin,
		Publisher:    publisher,
		Currency:     currency,
		Cashier:      cashier,
		Customer:     customer,
		Store:        store,
		Item:         item,
		PaymentTypes: paymentTypes,

		LongTermDBTx:            beginLongTermDBTx(longTermDatabase),
		ShardDBTx:               beginShardDBTx(shardingDatabase),
		GetShardIndexByDateTime: getShardIndexByDateTime(config.ShardingDatabase.Shards, config.Date),
		GetShardWhereQuery:      getShardWhereQuery(config.ShardingDatabase.Shards, config.Date),
	}, nil
}
