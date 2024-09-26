package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "172.105.54.158"
	port     = 5433
	user     = "positron"
	password = "B6W4qLAx80LrMX"
	dbname   = "positron"
)

type DB struct {
	*sql.DB
}
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func NewSqlDB(db_conf DBConfig) *DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db_conf.Host, db_conf.Port, db_conf.User, db_conf.Password, db_conf.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return &DB{
		DB: db,
	}
}
