package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"idea_server/util"
	"sync"
	"time"
)

var (
	db   *sql.DB
	once sync.Once
)

func getInstance() *sql.DB {
	if db == nil {
		once.Do(func() {
			cfg := util.LoadMysqlCfg()
			var err error
			db, err = sql.Open("mysql", cfg.Source)
			if err != nil {
				panic(err)
			}
			db.SetConnMaxLifetime(time.Minute * 3)
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(10)
		})
	}
	return db
}
