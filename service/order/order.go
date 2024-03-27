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
	"skripsi-be/repository/scheduler"
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

	FindInsightBasedOnInterval(ctx context.Context, interval string, offset int) (result.GetAggregateOrderService, error)
}

type Config struct {
	IsUsingSharding   bool
	Shards            config.Shards
	Date              config.Date
	KappaArchitecture config.KappaArchitecture
}

type service struct {
	config        Config
	orderRepo     order.Repository
	publisherRepo publisher.Repository
	currencyRepo  currency.Repository
	cashierRepo   cashier.Repository
	customerRepo  customer.Repository
	storeRepo     store.Repository
	itemRepo      item.Repository
	scheduler     scheduler.Repository

	beginShardTx            func(ctx context.Context, dbIdx int) (*sqlx.Tx, error)
	beginLongTermTx         func(ctx context.Context) (*sqlx.Tx, error)
	beginGeneralTx          func(ctx context.Context) (*sqlx.Tx, error)
	getShardIndexByDateTime func(date time.Time) (int, error)
	getShardWhereQuery      func(startDate time.Time, endDate time.Time) ([]result.ShardTimeSeriesWhereQuery, error)
}

func New(
	config Config,
	orderRepo order.Repository,
	publisherRepo publisher.Repository,
	currencyRepo currency.Repository,
	cashierRepo cashier.Repository,
	customerRepo customer.Repository,
	storeRepo store.Repository,
	itemRepo item.Repository,
	scheduler scheduler.Repository,

	beginShardTx func(ctx context.Context, dbIdx int) (*sqlx.Tx, error),
	beginLongTermTx func(ctx context.Context) (*sqlx.Tx, error),
	beginGeneralTx func(ctx context.Context) (*sqlx.Tx, error),
	getShardIndexByDateTime func(date time.Time) (int, error),
	getShardWhereQuery func(startDate time.Time, endDate time.Time) ([]result.ShardTimeSeriesWhereQuery, error),
) Service {
	return service{
		config:        config,
		orderRepo:     orderRepo,
		publisherRepo: publisherRepo,
		currencyRepo:  currencyRepo,
		cashierRepo:   cashierRepo,
		customerRepo:  customerRepo,
		storeRepo:     storeRepo,
		itemRepo:      itemRepo,
		scheduler:     scheduler,

		beginShardTx:            beginShardTx,
		beginLongTermTx:         beginLongTermTx,
		beginGeneralTx:          beginGeneralTx,
		getShardIndexByDateTime: getShardIndexByDateTime,
		getShardWhereQuery:      getShardWhereQuery,
	}
}
