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

// mysql 务必预编译sql语句
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

func Close() {
	if db != nil {
		_ = db.Close()
	}
}

func IsExistEmail(mail string) (bool, error) {
	db := getInstance()
	stmt, _ := db.Prepare("select id from users where email = ?")
	rows, err := stmt.Query(mail)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func Register(email, passwd string) {

}
