package routes

import (
	"github.com/gofiber/fiber/v2"
)

type ServerImpl interface {
}

type Server struct {
	*fiber.App
}

func NewApiServer() *fiber.App {
	return fiber.New()
}

// initialize routes
func (s *Server) InitializeRoutes(app *fiber.App) {

}
