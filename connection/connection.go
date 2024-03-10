package connection

import (
	"skripsi-be/config"

	pkgSQLX "skripsi-be/pkg/db/sqlx"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jmoiron/sqlx"
)

type Connection struct {
	GeneralDatabase  *sqlx.DB
	LongTermDatabase *sqlx.DB
	ShardingDatabase []*sqlx.DB
	KafkaProducer    *kafka.Producer
}

func New(config config.Config) (Connection, error) {
	generalDatabase, err := pkgSQLX.InitSQLX(config.GeneralDatabase.URIConnection)
	if err != nil {
		return Connection{}, err
	}

	longTermDatabase, err := pkgSQLX.InitSQLX(config.LongTermDatabase.URIConnection)
	if err != nil {
		return Connection{}, err
	}

	var shardingDatabase []*sqlx.DB
	for _, sharding := range config.ShardingDatabase.Shards {
		shardingDB, err := pkgSQLX.InitSQLX(sharding.URIConnection)
		if err != nil {
			return Connection{}, err
		}

		shardingDatabase = append(shardingDatabase, shardingDB)
	}

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.Server,
		"acks":              "all",
	})
	if err != nil {
		return Connection{}, err
	}

	return Connection{
		LongTermDatabase: longTermDatabase,
		ShardingDatabase: shardingDatabase,
		GeneralDatabase:  generalDatabase,
		KafkaProducer:    kafkaProducer,
	}, nil
}

func (c *Connection) Close() error {
	err := c.GeneralDatabase.Close()
	if err != nil {
		return err
	}

	err = c.LongTermDatabase.Close()
	if err != nil {
		return err
	}

	for idx := range c.ShardingDatabase {
		err = c.ShardingDatabase[idx].Close()
		if err != nil {
			return err
		}
	}

	c.KafkaProducer.Flush(10 * 1000)
	c.KafkaProducer.Close()

	return nil
}
