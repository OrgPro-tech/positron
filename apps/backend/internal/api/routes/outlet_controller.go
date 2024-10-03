package routes

import (
	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/OrgPro-tech/positron/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) CreateOutlet(c *fiber.Ctx) error {

	type outletData struct {
		OutletName    string `json:"outlet_name" validate:"required"`
		OutletAddress string `json:"outlet_address" validate:"required"`
		OutletPin     int32  `json:"outlet_pin" validate:"required"`
		OutletCity    string `json:"outlet_city" validate:"required"`
		OutletState   string `json:"outlet_state" validate:"required"`
		OutletCountry string `json:"outlet_country" validate:"required"`
		BusinessID    int32  `json:"business_id" validate:"required"`
	}
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
