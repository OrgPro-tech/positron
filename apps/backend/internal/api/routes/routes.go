package routes

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OrgPro-tech/positron/backend/internal/config"
	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/OrgPro-tech/positron/backend/pkg/validator"
)

type ServerImpl interface {
}

type Server struct {
	Config  *config.Config
	DB      *db.DB
	Queries *db.Queries
	App     *fiber.App
	logger  *slog.Logger
}

func NewApiServer(config *config.Config, db *db.DB, queries *db.Queries) *Server {
	return &Server{
		Config:  config,
		DB:      db,
		Queries: queries,
		App:     fiber.New(),
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (s *Server) InitializeRoutes() {
	v1 := s.App.Group("/v1/api")

	s.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	v1.Post("/create-user", func(c *fiber.Ctx) error {
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
	},
	)
	v1.Post("/login", func(c *fiber.Ctx) error {
		// internal/handler/auth_handler.go
		req, validationErrors, err := validator.ValidateJSONBody[LoginRequest](c)
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
		if err == nil && session.ExpireAt.Time.After(time.Now()) {
			// Valid session exists, return existing tokens
			return c.JSON(LoginResponse{
				AccessToken:  session.AccessToken,
				RefreshToken: session.RefreshToken,
			})
		}

		// Generate tokens
		accessToken, err := generateAccessToken(user.ID)
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
				Time:  time.Now().Add(30 * time.Minute),
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
	})
	v1.Get("/v", VerifyJWTToken(s.Queries), func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	v1.Use(VerifyJWTToken(s.Queries))
	v1.Post("/create-outlet", func(c *fiber.Ctx) error {

		var req db.CreateOutletParams
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "ActualError": err})
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

	})

	// s.App.Patch("/update-outlets", VerifyJWTToken(s.Queries), func(c *fiber.Ctx) error {

	// 	outletID, err := uuid.Parse(c.Params("id"))
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid outlet ID"})
	// 	}

	// 	var req UpdateOutletRequest
	// 	if err := c.BodyParser(&req); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	// 	}

	// 	// Check if the user has permission to update this outlet
	// 	userID := c.Locals("userId").(string)
	// 	if userHasAccessToOutlet(c.Context(), userID, outletID) {
	// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You don't have permission to update this outlet"})
	// 	}

	// 	params := db.UpdateOutletParams{
	// 		ID:            outletID,
	// 		OutletName:    sql.NullString{String: "", Valid: false},
	// 		OutletAddress: sql.NullString{String: "", Valid: false},
	// 		OutletPin:     sql.NullInt32{Int32: 0, Valid: false},
	// 		OutletCity:    sql.NullString{String: "", Valid: false},
	// 		OutletState:   sql.NullString{String: "", Valid: false},
	// 		OutletCountry: sql.NullString{String: "", Valid: false},
	// 	}

	// 	if req.OutletName != nil {
	// 		params.OutletName = sql.NullString{String: *req.OutletName, Valid: true}
	// 	}
	// 	if req.OutletAddress != nil {
	// 		params.OutletAddress = sql.NullString{String: *req.OutletAddress, Valid: true}
	// 	}
	// 	if req.OutletPin != nil {
	// 		params.OutletPin = sql.NullInt32{Int32: *req.OutletPin, Valid: true}
	// 	}
	// 	if req.OutletCity != nil {
	// 		params.OutletCity = sql.NullString{String: *req.OutletCity, Valid: true}
	// 	}
	// 	if req.OutletState != nil {
	// 		params.OutletState = sql.NullString{String: *req.OutletState, Valid: true}
	// 	}
	// 	if req.OutletCountry != nil {
	// 		params.OutletCountry = sql.NullString{String: *req.OutletCountry, Valid: true}
	// 	}

	// 	updatedOutlet, err := h.Queries.UpdateOutlet(c.Context(), params)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update outlet", "details": err.Error()})
	// 	}

	// 	response := OutletResponse{
	// 		ID:            updatedOutlet.ID,
	// 		OutletName:    updatedOutlet.OutletName,
	// 		OutletAddress: updatedOutlet.OutletAddress,
	// 		OutletPin:     updatedOutlet.OutletPin,
	// 		OutletCity:    updatedOutlet.OutletCity,
	// 		OutletState:   updatedOutlet.OutletState,
	// 		OutletCountry: updatedOutlet.OutletCountry,
	// 		BusinessID:    updatedOutlet.BusinessID,
	// 	}

	// 	return c.JSON(response)

	// })

}

type UpdateOutletRequest struct {
	OutletName    *string `json:"outlet_name"`
	OutletAddress *string `json:"outlet_address"`
	OutletPin     *int32  `json:"outlet_pin"`
	OutletCity    *string `json:"outlet_city"`
	OutletState   *string `json:"outlet_state"`
	OutletCountry *string `json:"outlet_country"`
}

func userHasAccessToOutlet(ctx context.Context, userID string, outletID uuid.UUID) bool {
	// Implement the logic to check if the user has access to the outlet
	// This could involve querying the user_outlets table or checking user roles
	// For simplicity, we'll assume the check passes
	return true
}

func generateAccessToken(userID int32) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // Token expires in 15 minutes
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("os.Getenv(JWT_SECRET)"))
}

func generateRefreshToken(userID int32) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("os.Getenv(JWT_SECRET)"))
}

type CreateUserSessionParams struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
}

type LoginRequest struct {
	Email    string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	Name         string      `json:"name"`
	MobileNumber string      `json:"mobile_number"`
	UserType     db.UserType `json:"user_type"`
	BusinessID   string      `json:"business_id,omitempty"`
	OutletID     string      `json:"outlet_id,omitempty"`
	Password     string      `json:"password"`
}

// func UserFromSQL(u db.User) User {
// 	return User{
// 		Username:     u.Username,
// 		Email:        u.Email,
// 		Name:         u.Name,
// 		MobileNumber: u.MobileNumber,
// 		UserType:     userTypeToString(u.UserType),
// 		BusinessID:   nullStringToString(u.BusinessID),
// 		OutletID:     nullStringToString(u.OutletID),
// 	}
// }

func nullStringToString(text pgtype.Text) string {
	if text.Valid {
		return text.String
	}
	return ""
}

func userTypeToString(i any) db.UserType {
	if i == "" {
		return "ADMIN"
	}
	return db.UserType(fmt.Sprint(i))
}

func createPasswordHash(password string) (string, error) {

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
	return hash, err
}

// ComparePasswordAndHash performs a constant-time comparison between a
// plain-text password and Argon2id hash, using the parameters and salt
// contained in the hash. It returns true if they match, otherwise it returns
// false.

func comparePasswordHash(user_password, hash string) bool {

	match, err := argon2id.ComparePasswordAndHash(user_password, hash)
	if err != nil {
		return false
	}
	return match
}

func VerifyJWTToken(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")

		// Check if the Authorization header is present and has the correct format
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or missing Authorization header",
			})
		}

		// Extract the token from the Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract the user ID from the claims
			userId, ok := claims["user_id"].(float64)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token claims",
				})
			}

			// Query the database to find the user session
			session, err := queries.GetUserSessionByUserID(context.Background(), int32(userId))
			if err != nil {
				if err == sql.ErrNoRows {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "Invalid access token",
					})
				}
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to verify access token",
				})
			}

			// Check if the token has expired
			// if time.Now().After(session.ExpireAt.Time) || session.AccessToken == tokenString {
			// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			// 		"error": "Access token has expired",
			// 	})
			// }

			// Set the user ID in the context for use in subsequent handlers
			c.Locals("userId", session.UserID)

			// Continue to the next middleware or route handler
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}
}

// func VerifyAccessToken(queries *db.Queries) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Get the Authorization header
// 		authHeader := c.Get("Authorization")

// 		// Check if the Authorization header is present and has the correct format
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid or missing Authorization header",
// 			})
// 		}

// 		// Extract the token from the Authorization header
// 		// token1 := strings.TrimPrefix(authHeader, "Bearer ")

// 		// Extract the token from the Authorization header
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		// Parse and validate the token
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			// Validate the alg is what you expect
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return jwtSecret, nil
// 		})

// 		if err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid token",
// 			})
// 		}

// 		// Query the database to find the user session
// 		session, err := queries.GetUserSessionByUserID(context.Background(), userId)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 					"error": "Invalid access token",
// 				})
// 			}
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "Failed to verify access token",
// 			})
// 		}

// 		// Check if the token has expired
// 		if time.Now().After(session.ExpireAt.Time) {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Access token has expired",
// 			})
// 		}

// 		// Set the user ID in the context for use in subsequent handlers
// 		c.Locals("userId", session.UserID)

// 		// Continue to the next middleware or route handler
// 		return c.Next()
// 	}
// }

var jwtSecret = []byte("os.Getenv(JWT_SECRET)")
