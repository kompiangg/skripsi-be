package sqlx

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func InitSQLX(dsn string) (db *sqlx.DB, err error) {
	if dsn == "" {
		return db, errors.New("[ERROR] config must not be nil")
	}

	delayTime := 1
	maxRetry := 5
	currRetry := 0

	for {
		db, err = sqlx.Connect("postgres", dsn)
		if err != nil && delayTime <= 20 {
			log.Error().Msgf("failed on creating connection on database, try to reconnect in %d second, err: %v", delayTime, err.Error())

			if currRetry == maxRetry {
				log.Fatal().Msgf("failed to connect to database, retry: %d", currRetry)
				return nil, err
			}

			time.Sleep(time.Second * time.Duration(delayTime))
			delayTime++
			currRetry++

			continue
		} else if err == nil {
			break
		}

		return db, err
	}

	db.SetConnMaxIdleTime(time.Second)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	return db, nil
}
