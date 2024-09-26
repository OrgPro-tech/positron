package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

type DB struct {
	*pgx.Conn
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func Connect(config DBConfig) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, psqlInfo)

	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)

	if err != nil {
		return nil, err
	}

	return &DB{conn}, nil
}
