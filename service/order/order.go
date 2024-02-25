package order

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/repository/cashier"
	"skripsi-be/repository/currency"
	"skripsi-be/repository/customer"
	"skripsi-be/repository/item"
	"skripsi-be/repository/order"
	"skripsi-be/repository/publisher"
	"skripsi-be/repository/shard"
	"skripsi-be/repository/store"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
	"time"

	"github.com/jmoiron/sqlx"
)

type Service interface {
	MoveDataThroughShard(ctx context.Context) error
	InsertToShard(ctx context.Context, param params.ServiceInsertOrdersToShardParam) error
	InsertToLongTerm(ctx context.Context, param params.ServiceInsertOrdersToLongTermParam) error
	InsertToShardSeeder(ctx context.Context, param params.ServiceInsertOrdersToShardParam) error
	InsertToLongTermSeeder(ctx context.Context, param params.ServiceInsertOrdersToLongTermParam) error
	FindOrder(ctx context.Context, param params.FindOrderService) (allOrders []result.Order, err error)
	IngestOrder(ctx context.Context, param []params.ServiceIngestionOrder) ([]result.ServiceIngestOrder, error)
	TransformOrder(ctx context.Context, param []params.ServiceTransformOrder) error

	FindBriefInformationOrder(ctx context.Context, param params.FindOrderService) (orders []result.OrderBriefInformation, err error)
	FindOrderDetails(ctx context.Context, param params.FindOrderDetailsService) (res result.Order, err error)
}

type Config struct {
	IsUsingSharding bool
	Shards          config.Shards
	Date            config.Date
}

type service struct {
	config        Config
	shardRepo     shard.Repository
	orderRepo     order.Repository
	publisherRepo publisher.Repository
	currencyRepo  currency.Repository
	cashierRepo   cashier.Repository
	customerRepo  customer.Repository
	storeRepo     store.Repository
	itemRepo      item.Repository

	beginShardTx            func(ctx context.Context, dbIdx int) (*sqlx.Tx, error)
	beginLongTermTx         func(ctx context.Context) (*sqlx.Tx, error)
	getShardIndexByDateTime func(date time.Time) (int, error)
	getShardWhereQuery      func(startDate time.Time, endDate time.Time) ([]result.ShardTimeSeriesWhereQuery, error)
}

func New(
	config Config,
	shardRepo shard.Repository,
	orderRepo order.Repository,
	publisherRepo publisher.Repository,
	currencyRepo currency.Repository,
	cashierRepo cashier.Repository,
	customerRepo customer.Repository,
	storeRepo store.Repository,
	itemRepo item.Repository,

	beginShardTx func(ctx context.Context, dbIdx int) (*sqlx.Tx, error),
	beginLongTermTx func(ctx context.Context) (*sqlx.Tx, error),
	getShardIndexByDateTime func(date time.Time) (int, error),
	getShardWhereQuery func(startDate time.Time, endDate time.Time) ([]result.ShardTimeSeriesWhereQuery, error),
) Service {
	return service{
		config:        config,
		shardRepo:     shardRepo,
		orderRepo:     orderRepo,
		publisherRepo: publisherRepo,
		currencyRepo:  currencyRepo,
		cashierRepo:   cashierRepo,
		customerRepo:  customerRepo,
		storeRepo:     storeRepo,
		itemRepo:      itemRepo,

		beginShardTx:            beginShardTx,
		beginLongTermTx:         beginLongTermTx,
		getShardIndexByDateTime: getShardIndexByDateTime,
		getShardWhereQuery:      getShardWhereQuery,
	}
}
