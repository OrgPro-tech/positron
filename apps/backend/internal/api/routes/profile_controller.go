package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetProfile(c *fiber.Ctx) error {

	userID := c.Locals("userId").(int32)

	profile, err := s.Queries.GetUserProfile(c.Context(), (userID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user profile"})
	}

	return c.JSON(profile)
}
