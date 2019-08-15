package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SQLDB ...
//
var SQLDB *sql.DB

func init() {
	log.Println("sqllite driver init")
	//打开数据库，如果不存在，则创建
	var err error
	SQLDB, err = sql.Open("sqlite3", "./rec.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	//SQLDB.SetMaxIdleConns(20)
	//SQLDB.SetMaxOpenConns(20)
	err = SQLDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("sqllite driver init ok!")

}
