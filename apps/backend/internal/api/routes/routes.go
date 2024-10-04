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

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OrgPro-tech/positron/backend/internal/config"
	"github.com/OrgPro-tech/positron/backend/internal/db"
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

	v1.Post("/create-user", s.CreateUser)
	v1.Post("/login", s.Login)
	v1.Get("/v", VerifyJWTToken(s.Queries), func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	v1.Use(VerifyJWTToken(s.Queries))
	v1.Post("/create-outlet", s.CreateOutlet)
	v1.Get("/profile", s.GetProfile)
	v1.Post("/create-category", s.CreateCategory)
	v1.Get("/get-category", s.GetAllCategories)
	v1.Post("/create-menu", s.CreateMenuItem)

}
func verifyRefreshToken(tokenString string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
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

func generateAccessToken(userID, businessId int32) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID,
		"business_id": businessId,
		"exp":         time.Now().Add(24 * time.Hour).Unix(), // Token expires in 15 minutes
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
				"error":   "Invalid token",
				"message": err.Error(),
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

			businessId, ok := claims["business_id"].(float64)
			if !ok && businessId == 0 {
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
			c.Locals("business_id", int32(businessId))

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

func SendErrResponse(ctx *fiber.Ctx, err error, statusCode int) error {
	return ctx.Status(statusCode).JSON(&fiber.Map{
		"status": statusCode,
		"error": &fiber.Map{
			"message": err.Error(),
		},
	})
}

func SendSuccessResponse(ctx *fiber.Ctx, message string, data interface{}, statusCode int) error {
	return ctx.Status(statusCode).JSON(&fiber.Map{
		"status":  statusCode,
		"message": message,
		"data":    data,
	})
}
