package db

// import (
// 	"database/sql"
// 	"fmt"

// 	"github.com/jackc/pgx/v5"
// 	_ "github.com/lib/pq"
// )

// type DB struct {
// 	*sql.DB
// }
// type DBConfig struct {
// 	Host     string
// 	Port     string
// 	User     string
// 	Password string
// 	Dbname   string
// }

// func Connect(db_conf DBConfig) *pgx.DB {
// 	// psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
// 	// 	"password=%s dbname=%s sslmode=disable",
// 	// 	db_conf.Host, db_conf.Port, db_conf.User, db_conf.Password, db_conf.Dbname)
// 	// db, err := sql.Open("postgres", psqlInfo)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// defer db.Close()

// 	conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close(ctx)

// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Successfully connected!")
// 	return &DB{
// 		DB: db,
// 	}
// }
