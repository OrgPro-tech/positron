package routes

import (
	"context"
	"time"

	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/OrgPro-tech/positron/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) Login(c *fiber.Ctx) error {
	// internal/handler/auth_handler.go
	req, validationErrors, err := validator.ValidateJSONBody[LoginRequest](c)
	if err != nil {
		return SendErrResponse(c, err, fiber.StatusBadRequest)
		// 	c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"error": err.Error(),
		// 	})
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	// Get user from database
	user, err := s.Queries.GetUserByUsernameOrEmail(c.Context(), req.Email) //.GetUserByUsername(context.Background(), req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials", "error_cred": err})
	}

	// Check password
	if !comparePasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials", "error_password_check": err})
	}

	// Check if valid session exists
	session, err := s.Queries.GetUserSessionByUserID(c.Context(), user.ID)
	if err == nil && session.ExpireAt.Time.Before(time.Now()) {
		// Valid session exists, return existing tokens
		return c.JSON(LoginResponse{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
		})
	}

	// Generate tokens
	accessToken, err := generateAccessToken(user.ID, user.BusinessID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
	}

	// Start a transaction
	tx, err := s.DB.Begin(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}
	defer tx.Rollback(c.Context())

	qtx := s.Queries.WithTx(tx)

	// Delete existing sessions for the user
	if err := qtx.DeleteUserSessions(c.Context(), user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete old sessions"})
	}

	// Create new session
	newSession, err := qtx.CreateUserSession(context.Background(), db.CreateUserSessionParams{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt: pgtype.Timestamp{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		},
	},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
	}

	// Commit the transaction
	if err := tx.Commit(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	// // Save refresh token to database
	// sessiondata, err := s.Queries.CreateUserSession(context.Background(), db.CreateUserSessionParams{
	// 	UserID:       user.ID,
	// 	AccessToken:  accessToken,
	// 	RefreshToken: refreshToken,
	// 	ExpireAt: pgtype.Timestamp{
	// 		Time:  time.Now().Add(30 * time.Minute),
	// 		Valid: true,
	// 	},
	// },
	// )
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session", "errorData": err})
	// }

	return c.JSON(LoginResponse{
		AccessToken:  newSession.AccessToken,
		RefreshToken: newSession.RefreshToken,
	})
}

type CreateUserWithBusinessParams struct {
	User struct {
		Username     string `json:"username" validate:"required"`
		Password     string `json:"password" validate:"required"`
		Email        string `json:"email" validate:"required"`
		Name         string `json:"name" validate:"required"`
		MobileNumber string `json:"mobile_number" validate:"required"`
		UserType     string `json:"user_type" validate:"required"`
	} `json:"user" validate:"required"`
	Business struct {
		ContactPersonName         string `json:"contact_person_name" validate:"required"`
		ContactPersonEmail        string `json:"contact_person_email" validate:"required"`
		ContactPersonMobileNumber string `json:"contact_person_mobile_number" validate:"required"`
		CompanyName               string `json:"company_name" validate:"required"`
		Address                   string `json:"address" validate:"required"`
		Pin                       int32  `json:"pin" validate:"required"`
		City                      string `json:"city" validate:"required"`
		State                     string `json:"state" validate:"required"`
		Country                   string `json:"country" validate:"required"`
		BusinessType              string `json:"business_type" validate:"required"`
		Gst                       string `json:"gst" validate:"required"`
		Pan                       string `json:"pan" validate:"required"`
		BankAccountNumber         string `json:"bank_account_number" validate:"required"`
		BankName                  string `json:"bank_name" validate:"required"`
		IfscCode                  string `json:"ifsc_code" validate:"required"`
		AccountType               string `json:"account_type" validate:"required"`
		AccountHolderName         string `json:"account_holder_name" validate:"required"`
	} `json:"Business" validate:"required"`
}

func (s *Server) CreateUser(c *fiber.Ctx) error {

	reqBody, validationErrors, err := validator.ValidateJSONBody[CreateUserWithBusinessParams](c)
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
	// var reqBody CreateUserWithBusinessParams
	// if err := c.BodyParser(&reqBody); err != nil {
	// 	s.logger.Error("Failed to parse request body", "error", err)
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Invalid request body",
	// 	})
	// }
	req := db.CreateUserWithBusinessParams{
		ContactPersonName:         reqBody.Business.ContactPersonName,
		ContactPersonEmail:        reqBody.Business.ContactPersonEmail,
		ContactPersonMobileNumber: reqBody.Business.ContactPersonMobileNumber,
		CompanyName:               reqBody.Business.CompanyName,
		Address:                   reqBody.Business.Address,
		Pin:                       reqBody.Business.Pin,
		City:                      reqBody.Business.City,
		State:                     reqBody.Business.State,
		Country:                   reqBody.Business.Country,
		BusinessType:              reqBody.Business.BusinessType,
		Gst:                       reqBody.Business.Gst,
		Pan:                       reqBody.Business.Pan,
		BankAccountNumber:         reqBody.Business.BankAccountNumber,
		BankName:                  reqBody.Business.BankName,
		IfscCode:                  reqBody.Business.IfscCode,
		AccountType:               reqBody.Business.AccountType,
		AccountHolderName:         reqBody.Business.AccountHolderName,
		Username:                  reqBody.User.Username,
		Password:                  reqBody.User.Password,
		Email:                     reqBody.User.Email,
		Name:                      reqBody.User.Name,
		MobileNumber:              reqBody.User.MobileNumber,
		UserType:                  userTypeToString(reqBody.User.UserType),
	}

	// Hash the password before storing
	hashedPassword, err := createPasswordHash(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process user data",
		})
	}
	req.Password = hashedPassword

	user, err := s.Queries.CreateUserWithBusiness(c.Context(), req)
	if err != nil {
		s.logger.Error("Failed to create user and business", "error", err)
		// if db.IsUniqueViolation(err) {
		// 	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		// 		"error": "User with this username or email already exists",
		// 	})
		// }
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user and business",
		})
	}

	s.logger.Info("User and business created successfully", "user_id", user.ID)
	return c.Status(fiber.StatusCreated).JSON(user)
	// return c.Status(fiber.StatusCreated).JSON((createdUser))
}
