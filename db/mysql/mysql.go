package mysql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"idea_server/util"
	"log"
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
	defer stmt.Close()
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

func Register(username, email, passwd string) error {
	db := getInstance()
	stmt, _ := db.Prepare("insert into users (username, email, passwd) VALUES (?,?,?)")
	defer stmt.Close()
	_, err := stmt.Exec(username, email, passwd)
	if err != nil {
		return err
	}
	return nil
}

func Login(email, passwd string) bool {
	db := getInstance()
	stmt, _ := db.Prepare("select id from users where email = ? and passwd = ?")
	defer stmt.Close()
	rows, err := stmt.Query(email, passwd)
	if err != nil {
		log.Println(err)
		return false
	}
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}

func GetID(email string) (string, error) {
	db := getInstance()
	stmt, _ := db.Prepare("select id from users where email = ?")
	defer stmt.Close()
	rows, err := stmt.Query(email)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return "", err
		}
		return id, nil
	}
	return "", errors.New("does not exist")
}

// AllUserGroup Return to the group chat where the user is
// return nil on error
func AllUserGroup(id string) []string {
	db := getInstance()
	stmt, _ := db.Prepare("select gid from member where uid = ?")
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	ans := make([]string, 0)
	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err != nil {
			return nil
		}
		ans = append(ans, str)
	}
	return ans
}
