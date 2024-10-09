package routes

import (
	"encoding/json"
	"errors"
	"math/big"

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

type CreateMenuItemRequest struct {
	CategoryID    int32           `json:"category_id"`
	Name          string          `json:"name" validate:"required,max=100"`
	Description   *string         `json:"description"`
	Price         float32         `json:"price" validate:"required"`
	IsVegetarian  bool            `json:"is_vegetarian"`
	SpiceLevel    *string         `json:"spice_level" validate:"omitempty,oneof=Mild Medium Hot ExtraHot"`
	IsAvailable   bool            `json:"is_available"`
	BusinessID    int32           `json:"business_id" validate:"required"`
	Code          string          `json:"code" validate:"required"`
	TaxPercentage int             `json:"tax_percentage" validate:"required,min=0,max=100"`
	SizeType      string          `json:"size_type" validate:"required,oneof=GRAM PIECE"`
	Variation     json.RawMessage `json:"variation" validate:"omitempty,json"`
	Customizable  bool            `json:"customizable"`
	Image         string          `json:"image" validate:"omitempty,url"`
}

func float32ToPgNumeric(f float32) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   new(big.Int).SetInt64(int64(f * 100)), // Multiply by 100 to preserve 2 decimal places
		Exp:   -2,                                    // Set the exponent to -2 to account for the multiplication
		Valid: true,
	}
}
func (s *Server) CreateMenuItem(c *fiber.Ctx) error {
	businessID := c.Locals("business_id").(int32)

	if businessID == 0 {
		return SendErrResponse(c, errors.New("Invalid business ID"), fiber.StatusBadRequest)
	}

	req, validationErrors, err := validator.ValidateJSONBody[CreateMenuItemRequest](c)
	if err != nil {
		return SendErrResponse(c, err, fiber.StatusInternalServerError)
	}
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	// Create the menu item
	menuItem, err := s.Queries.CreateMenuItem(c.Context(), db.CreateMenuItemParams{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Description:   pgtype.Text{String: *req.Description, Valid: req.Description != nil},
		Price:         float32ToPgNumeric(req.Price),
		IsVegetarian:  req.IsVegetarian,
		SpiceLevel:    db.NullSpiceLevel{SpiceLevel: db.SpiceLevel(*req.SpiceLevel), Valid: req.SpiceLevel != nil},
		IsAvailable:   req.IsAvailable,
		BusinessID:    int32(businessID),
		Code:          req.Code,
		TaxPercentage: int32(req.TaxPercentage),
		SizeType:      db.SizeType(req.SizeType),
		Variation:     req.Variation,
		Customizable:  req.Customizable,
		Image: pgtype.Text{
			String: req.Image,
			Valid:  req.Image != "",
		},
	})
	if err != nil {
		return SendErrResponse(c, err, fiber.StatusInternalServerError)
	}

	return SendSuccessResponse(c, "Menu created successfully", menuItem, fiber.StatusCreated) //c.Status(fiber.StatusCreated).JSON(menuItem)
}
func (s *Server) GetAllMenuItemsByBusinessID(c *fiber.Ctx) error {
	businessID := c.Locals("business_id").(int32)

	if businessID == 0 {
		return SendErrResponse(c, errors.New("Invalid business ID"), fiber.StatusBadRequest)
	}

	menuItems, err := s.Queries.GetMenuItemsByBusinessID(c.Context(), int32(businessID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch menu items"})
	}

	return SendSuccessResponse(c, "Fetch successful", menuItems, fiber.StatusOK)
	//c.JSON(menuItems)
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

	return SendSuccessResponse(c, "Fetch successful", categories, fiber.StatusOK)
}
