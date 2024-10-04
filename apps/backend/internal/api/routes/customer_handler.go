package routes

import (
	"database/sql"
	"errors"

	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/OrgPro-tech/positron/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCustomerRequest struct {
	PhoneNumber string  `json:"phone_number" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	WhatsApp    *bool   `json:"whatsapp"`
	Email       *string `json:"email"`
	Address     *string `json:"address"`
	OutletID    int32   `json:"outlet_id" validate:"required"`
}

type UpdateCustomerRequest struct {
	PhoneNumber string  `json:"phone_number" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	WhatsApp    *bool   `json:"whatsapp"  validate:"omitempty"`
	Email       *string `json:"email" validate:"omitempty"`
	Address     *string `json:"address"  validate:"omitempty"`
	OutletID    int32   `json:"outlet_id" validate:"required"`
}

func (s *Server) CreateCustomer(c *fiber.Ctx) error {
	businessID := c.Locals("business_id").(int32)

	if businessID == 0 {
		return SendErrResponse(c, errors.New("Invalid business ID"), fiber.StatusBadRequest)
	}
	req, validationErrors, err := validator.ValidateJSONBody[CreateCustomerRequest](c)
	if err != nil {
		return SendErrResponse(c, err, fiber.StatusInternalServerError)
	}
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}
	customer, err := s.Queries.CreateCustomer(c.Context(), db.CreateCustomerParams{
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Whatsapp: pgtype.Bool{
			Bool:  *req.WhatsApp,
			Valid: req.WhatsApp != nil,
		},
		// WhatsApp:    sql.NullBool{Bool: *req.WhatsApp, Valid: req.WhatsApp != nil},
		Email: pgtype.Text{
			String: *req.Email,
			Valid:  req.Email != nil,
		}, //sql.NullString{String: *req.Email, Valid: req.Email != nil},
		Address: pgtype.Text{
			String: *req.Address,
			Valid:  req.Address != nil,
		}, //sql.NullString{String: *req.Address, Valid: req.Address != nil},
		OutletID:   req.OutletID,
		BusinessID: businessID,
	})
	if err != nil {
		return SendErrResponse(c, err, fiber.StatusInternalServerError)
	}

	return SendSuccessResponse(c, "Customer created successfully", customer, fiber.StatusCreated) //c.Status(fiber.StatusCreated).JSON(menuItem)
}

func (s *Server) GetCustomer(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid customer ID"})
	}

	customer, err := s.Queries.GetCustomerByID(c.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch customer"})
	}

	return c.JSON(customer)
}

func (s *Server) GetCustomersByBusinessId(c *fiber.Ctx) error {
	businessID := c.Locals("business_id").(int32)
	if businessID == 0 {
		return SendErrResponse(c, errors.New("Invalid business ID"), fiber.StatusBadRequest)
	}

	customers, err := s.Queries.GetCustomersByBusinessID(c.Context(), int32(businessID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch customers"})
	}

	return c.JSON(customers)
}
func (s *Server) GetCustomersByOutletId(c *fiber.Ctx) error {
	outlet_id, err := c.ParamsInt("outlet_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid business ID"})
	}

	customers, err := s.Queries.GetCustomersByOutletId(c.Context(), int32(outlet_id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch customers"})
	}

	return c.JSON(customers)
}

// func (s *Server) UpdateCustomer(c *fiber.Ctx) error {
// 	id, err := c.ParamsInt("id")
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid customer ID"})
// 	}

// 	var req UpdateCustomerRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
// 	}

// 	if err := c.Validate(req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	customer, err := h.Queries.UpdateCustomer(c.Context(), db.UpdateCustomerParams{
// 		ID:          int32(id),
// 		PhoneNumber: req.PhoneNumber,
// 		Name:        req.Name,
// 		WhatsApp:    sql.NullBool{Bool: *req.WhatsApp, Valid: req.WhatsApp != nil},
// 		Email:       sql.NullString{String: *req.Email, Valid: req.Email != nil},
// 		Address:     sql.NullString{String: *req.Address, Valid: req.Address != nil},
// 		OutletID:    req.OutletID,
// 	})
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update customer"})
// 	}

// 	return c.JSON(customer)
// }
