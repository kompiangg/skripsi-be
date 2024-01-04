package connection

import (
	"skripsi-be/config"

	pkgSQLX "skripsi-be/pkg/db/sqlx"
	pkgRedis "skripsi-be/pkg/redis"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Connection struct {
	GeneralDatabase  *sqlx.DB
	LongTermDatabase *sqlx.DB
	ShardingDatabase []*sqlx.DB
	Redis            *redis.Client
}

func New(config config.Config) (Connection, error) {
	redis, err := pkgRedis.New(pkgRedis.RedisConfig{
		Hostname: config.Redis.Hostname,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if err != nil {
		return Connection{}, err
	}

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

	return Connection{
		Redis:            redis,
		LongTermDatabase: longTermDatabase,
		ShardingDatabase: shardingDatabase,
		GeneralDatabase:  generalDatabase,
	}, nil
}

func (c *Connection) Close() error {
	err := c.Redis.Close()
	if err != nil {
		return err
	}

	err = c.GeneralDatabase.Close()
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

	return nil
}
