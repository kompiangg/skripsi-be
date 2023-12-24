package sqlx

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitSQLX(dsn string) (db *sqlx.DB, err error) {
	if dsn == "" {
		return db, errors.New("[ERROR] config must not be nil")
	}

	delayTime := 1

	for {
		db, err = sqlx.Connect("postgres", dsn)
		if err != nil && delayTime <= 20 {
			fmt.Printf("error (%v): failed on creating connection on database, try to reconnect in %d second\n", err.Error(), delayTime)

			time.Sleep(time.Second * time.Duration(delayTime))
			delayTime++

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
