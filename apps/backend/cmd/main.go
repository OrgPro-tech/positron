package main

import (
	"fmt"
	"log"

	"github.com/OrgPro-tech/positron/backend/internal/api/routes"
	"github.com/OrgPro-tech/positron/backend/internal/config"
	"github.com/OrgPro-tech/positron/backend/internal/db"
)

func main() {
	envConfig := config.NewConfig()
	fmt.Printf("envConfig: %v\n", envConfig)
	dbConn, err := db.Connect(db.DBConfig{
		Host:     envConfig.Host,
		Port:     envConfig.DB_Port,
		User:     envConfig.Username,
		Password: envConfig.Password,
		Dbname:   envConfig.Database_Name,
	})

	if err != nil {
		panic(err)
	}

	queries := db.New(dbConn)

	server := routes.NewApiServer(envConfig, dbConn, queries)
	server.InitializeRoutes()
	log.Fatal(server.App.Listen(":" + envConfig.ServerPort))

}
