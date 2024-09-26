package main

import (
	"log"

	"command-line-arguments/Users/shubhabanerjee/orgpro.tech/positron/apps/backend/internal/db/db.go"

	"github.com/OrgPro-tech/positron/backend/internal/api/routes"
	"github.com/OrgPro-tech/positron/backend/internal/config"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	config  *config.Config
	db      *db.DB
	app     *fiber.App
	queries db.Queries
}

func main() {

	envConfig := config.NewConfig()
	app := routes.NewApiServer()
	db := db.Connect(db.DBConfig{
		Host:     envConfig.Host,
		Port:     envConfig.DB_Port,
		User:     envConfig.Username,
		Password: envConfig.Password,
		Dbname:   envConfig.Database_Name,
	})
	queries := db.New(db)
	s := Server{
		config:  envConfig,
		db:      db,
		app:     app,
		queries: queries,
	}

	log.Fatal(s.app.Listen(":" + envConfig.ServerPort))

}
