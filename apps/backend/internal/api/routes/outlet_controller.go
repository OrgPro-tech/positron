package routes

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/OrgPro-tech/positron/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type outletData struct {
	OutletName    string `json:"outlet_name" validate:"required"`
	OutletAddress string `json:"outlet_address" validate:"required"`
	OutletPin     int32  `json:"outlet_pin" validate:"required"`
	OutletCity    string `json:"outlet_city" validate:"required"`
	OutletState   string `json:"outlet_state" validate:"required"`
	OutletCountry string `json:"outlet_country" validate:"required"`
	BusinessID    int32  `json:"business_id" validate:"required"`
}

func (s *Server) CreateOutlet(c *fiber.Ctx) error {

	req, validationErrors, err := validator.ValidateJSONBody[outletData](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	userID := c.Locals("userId").(int32)

	result, err := s.Queries.CreateOutletWithUserAssociation(c.Context(), db.CreateOutletWithUserAssociationParams{
		OutletName:    req.OutletName,
		OutletAddress: req.OutletAddress,
		OutletPin:     req.OutletPin,
		OutletCity:    req.OutletCity,
		OutletState:   req.OutletState,
		OutletCountry: req.OutletCountry,
		BusinessID:    (req.BusinessID),
		UserID:        (userID),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create outlet", "details": err.Error()})
	}

	response := db.CreateOutletWithUserAssociationRow{
		ID:            result.ID,
		OutletName:    result.OutletName,
		OutletAddress: result.OutletAddress,
		OutletPin:     result.OutletPin,
		OutletCity:    result.OutletCity,
		OutletState:   result.OutletState,
		OutletCountry: result.OutletCountry,
		BusinessID:    result.BusinessID,
		UserOutletID:  result.UserOutletID,
	}

	return c.Status(fiber.StatusCreated).JSON(response)

}

type CreateOutletMenuItemInput struct {
	OutletID    int32   `json:"outlet_id" validate:"required"`
	MenuItemID  int32   `json:"menu_item_id" validate:"required"`
	Price       float32 `json:"price" validate:"required,numeric"`
	IsAvailable bool    `json:"is_available"`
}

func (s *Server) AddOutletMenu(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int32)
	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user id",
		})
	}
	var input CreateOutletMenuItemInput
	input, validationErrors, err := validator.ValidateJSONBody[CreateOutletMenuItemInput](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}
	outletMenuItem, err := s.Queries.CreateOutletMenuItem(c.Context(), db.CreateOutletMenuItemParams{
		OutletID:    input.OutletID,
		MenuItemID:  input.MenuItemID,
		Price:       float32ToPgNumeric(input.Price),
		IsAvailable: input.IsAvailable,
		CreatedBy:   userID,
	})

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create outlet  menu item", "details": err.Error()})

	}

	return c.Status(fiber.StatusCreated).JSON(outletMenuItem)
}

func (s *Server) GetMenuByOutlet(c *fiber.Ctx) error {

	outletID, err := strconv.Atoi(c.Params("outletId"))
	if err != nil {
		return SendErrResponse(c, errors.New("Invalid outlet ID"), fiber.StatusInternalServerError)

	}

	outletMenuItems, err := s.Queries.ListOutletMenuItems(c.Context(), int32(outletID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list outlet menu items"})
	}

	return SendSuccessResponse(c, "Menu fetch successfully", outletMenuItems, fiber.StatusCreated)
}

func (s *Server) UpdateMenuByOutlet(c *fiber.Ctx) error {
	type UpdateOutletMenuItemInput struct {
		Price       *float32 `json:"price"`
		IsAvailable *bool    `json:"is_available"`
	}
	outletID, err := strconv.Atoi(c.Params("outletId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid outlet ID"})
	}

	menuItemID, err := strconv.Atoi(c.Params("menuId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu item ID"})
	}

	var input UpdateOutletMenuItemInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Fetch the current outlet menu item
	currentItem, err := s.Queries.GetOutletMenuItem(c.Context(), db.GetOutletMenuItemParams{
		OutletID:   int32(outletID),
		MenuItemID: int32(menuItemID),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Outlet menu item not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch outlet menu item"})
	}

	// Update only the fields that are provided
	updateParams := db.UpdateOutletMenuItemParams{
		OutletID:    int32(outletID),
		MenuItemID:  int32(menuItemID),
		Price:       currentItem.Price,
		IsAvailable: currentItem.IsAvailable,
	}

	if input.Price != nil {
		updateParams.Price = float32ToPgNumeric(*input.Price) //*input.Price
	}
	if input.IsAvailable != nil {
		updateParams.IsAvailable = *input.IsAvailable
	}

	updatedItem, err := s.Queries.UpdateOutletMenuItem(c.Context(), updateParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update outlet menu item"})
	}

	return SendSuccessResponse(c, "Menu fetch successfully", updatedItem, fiber.StatusCreated)
}
