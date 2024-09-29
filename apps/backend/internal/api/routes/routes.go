package routes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/OrgPro-tech/positron/backend/internal/config"
	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
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
	s.App.Post("/please-login", func(c *fiber.Ctx) error {
		var email = "vaibhav@itday.in"
		user, err := s.Queries.GetUser(context.Background(), email)

		if err != nil {
			return errors.New("invalid email id")
		}

		return c.JSON(map[string]string{
			"email":    user.Email,
			"password": user.Password,
		})
	})

	s.App.Post("/create-user", func(c *fiber.Ctx) error {
		// var user User
		// if err := c.BodyParser(&user); err != nil {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "errorData": err})
		// }
		// fmt.Printf("user: %v\n", user)

		// params := db.CreateUserParams{
		// 	Username:     user.Username,
		// 	Email:        user.Email,
		// 	Name:         user.Name,
		// 	MobileNumber: user.MobileNumber,
		// 	Password:     createPasswordHash(user.Password),
		// 	UserType:     userTypeToString(user.UserType),
		// }

		// createdUser, err := s.Queries.CreateUser(c.Context(), params)
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user", "actualError": err})
		// }
		var req db.CreateUserWithBusinessParams
		if err := c.BodyParser(&req); err != nil {
			s.logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// if err := validate.Struct(req); err != nil {
		// 	s.logger.Error("Input validation failed", "error", err)
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"error":   "Invalid input data",
		// 		"details": util.FormatValidationErrors(err),
		// 	})
		// }

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
	s.App.Post("/login", func(c *fiber.Ctx) error {
		// internal/handler/auth_handler.go

		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
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

		// Generate tokens
		accessToken, err := generateAccessToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
		}

		refreshToken, err := generateRefreshToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
		}

		// Save refresh token to database
		sessiondata, err := s.Queries.CreateUserSession(context.Background(), db.CreateUserSessionParams{
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session", "errorData": err})
		}

		return c.JSON(LoginResponse{
			AccessToken:  sessiondata.AccessToken,
			RefreshToken: sessiondata.RefreshToken,
		})
	})
	s.App.Get("/v", VerifyJWTToken(s.Queries), func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	s.App.Post("/create-outlet", VerifyJWTToken(s.Queries), func(c *fiber.Ctx) error {

		var req db.CreateOutletParams
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		userID := c.Locals("userId").(string)

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

}

func generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // Token expires in 15 minutes
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("os.Getenv(JWT_SECRET)"))
}

func generateRefreshToken(userID string) (string, error) {
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
	Email    string `json:"username"`
	Password string `json:"password"`
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

func UserFromSQL(u db.User) User {
	return User{
		Username:     u.Username,
		Email:        u.Email,
		Name:         u.Name,
		MobileNumber: u.MobileNumber,
		UserType:     userTypeToString(u.UserType),
		BusinessID:   nullStringToString(u.BusinessID),
		OutletID:     nullStringToString(u.OutletID),
	}
}

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
			userId, ok := claims["user_id"].(string)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token claims",
				})
			}

			// Query the database to find the user session
			session, err := queries.GetUserSessionByUserID(context.Background(), userId)
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
