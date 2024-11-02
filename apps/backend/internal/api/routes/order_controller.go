package routes

import (
	"strings"

	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/OrgPro-tech/positron/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateOrder struct {
	CustomerID  int32   `json:"customer_id"`
	PhoneNumber string  `json:"phone_number"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Address     string  `json:"address"`
	OrderID     string  `json:"order_id"`
	Status      string  `json:"status"`
	GstAmount   float32 `json:"gst_amount"`
	TotalAmount float32 `json:"total_amount"`
	NetAmount   float32 `json:"net_amount"`
}

func (s *Server) CreateOrder(c *fiber.Ctx) error {
	req, validationErrors, err := validator.ValidateJSONBody[CreateOrder](c)
	if err != nil {
		return SendErrResponse(c, err, fiber.StatusInternalServerError)
	}
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	//get the user id from the token
	r, err := s.Queries.CreateOrder(c.Context(), db.CreateOrderParams{
		CustomerID:  req.CustomerID,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Email: pgtype.Text{
			String: req.Email,
			Valid:  req.Email != "",
		}, //
		Address: pgtype.Text{
			String: req.Address,
			Valid:  req.Address != "",
		}, //

		OrderID: req.OrderID,
		Status:  db.OrderStatus(strings.ToUpper(req.Status)),
		// Status:   , //  req.Status,
		GstAmount:   float32ToPgNumeric(req.GstAmount),
		TotalAmount: float32ToPgNumeric(req.TotalAmount),
		NetAmount:   float32ToPgNumeric(req.NetAmount),
	})
	if err != nil {
		return err
	}
	return c.JSON(r)
}
