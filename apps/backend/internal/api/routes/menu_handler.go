package routes

import (
	"errors"

	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/OrgPro-tech/positron/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,max=50"`
	Description *string `json:"description"`
	BusinessID  int32   `json:"business_id" validate:"required"`
}

func (s *Server) CreateCategory(c *fiber.Ctx) error {

	req, validationErrors, err := validator.ValidateJSONBody[CreateCategoryRequest](c)
	if err != nil {
		return SendErrResponse(c, err, fiber.StatusInternalServerError)
	}
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	category, err := s.Queries.CreateCategory(c.Context(), db.CreateCategoryParams{
		Name: req.Name,
		Description: pgtype.Text{
			String: description,
			Valid:  req.Description != nil,
		},
		BusinessID: req.BusinessID,
	})
	if err != nil {
		return SendErrResponse(c, err, fiber.StatusInternalServerError) // c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create category"})
	}

	return SendSuccessResponse(c, "Category successfully created", category, fiber.StatusCreated) //c.Status(fiber.StatusCreated).JSON(category)
}

func (s *Server) GetAllCategories(c *fiber.Ctx) error {
	businessID := c.Locals("business_id").(int32)

	if businessID == 0 {
		return SendErrResponse(c, errors.New("Invalid business ID"), fiber.StatusBadRequest)
	}

	categories, err := s.Queries.GetAllCategories(c.Context(), int32(businessID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}

	return SendSuccessResponse(c, "Category successfully created", categories, fiber.StatusOK)
}
