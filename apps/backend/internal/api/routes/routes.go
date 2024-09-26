package routes

import (
	"context"
	"errors"

	"github.com/OrgPro-tech/positron/backend/internal/config"
	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/gofiber/fiber/v2"
)

type ServerImpl interface {
}

type Server struct {
	Config  *config.Config
	DB      *db.DB
	Queries *db.Queries
	App     *fiber.App
}

func NewApiServer(config *config.Config, db *db.DB, queries *db.Queries) *Server {
	return &Server{
		Config:  config,
		DB:      db,
		Queries: queries,
		App:     fiber.New(),
	}
}

func (s *Server) InitializeRoutes() {
	s.App.Post("/please-login", func(c *fiber.Ctx) error {
		var email = "vaibhav@itday.in"
		user, err := s.Queries.GetUser(context.Background(), email)

		if err != nil {
			return errors.New("invalid email id")
		}

		return c.JSON(map[string]string{
			"email":    user.Email,
			"password": user.Password,
		})
	})
}
