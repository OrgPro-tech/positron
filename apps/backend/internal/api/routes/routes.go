package routes

import (
	"context"
	"errors"
	"fmt"

	"github.com/OrgPro-tech/positron/backend/internal/config"
	"github.com/OrgPro-tech/positron/backend/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type ServerImpl interface {
}

type Server struct {
	Config  *config.Config
	DB      *db.DB
	Queries *db.Queries
	App     *fiber.App
}

func NewApiServer(config *config.Config, db *db.DB, queries *db.Queries) *Server {
	return &Server{
		Config:  config,
		DB:      db,
		Queries: queries,
		App:     fiber.New(),
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
		// var params db.CreateUserParams
		// u := User{}
		// if err := c.BodyParser(&u); err != nil {
		// 	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		// }
		// fmt.Printf("params: %+v\n", params)
		// fmt.Printf("u: %+v\n", u)
		// user, err := s.Queries.CreateUser(context.Background(), db.CreateUserParams{
		// 	ID:       u.ID,
		// 	Name:     u.Name,
		// 	Email:    u.Email,
		// 	Password: u.Password,
		// 	BusinessId: pgtype.Text{
		// 		String: "cd1",
		// 		Valid:  true,
		// 	},
		// 	MobileNumber: u.MobileNumber,
		// 	OutletId:     params.OutletId, //u.OutletId,
		// 	UserType:     u.UserType,
		// 	Username:     u.Username,
		// })
		// if err != nil {
		// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user", "ActualError": err})
		// }

		// return c.Status(http.StatusCreated).JSON(user)
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		fmt.Printf("user: %v\n", user)

		params := db.CreateUserParams{
			Username:     user.Username,
			Email:        user.Email,
			Name:         user.Name,
			MobileNumber: user.MobileNumber,
			Password:     user.Password,
			UserType:     userTypeToString(user.UserType),
		}

		createdUser, err := s.Queries.CreateUser(c.Context(), params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user", "actualError": err})
		}

		return c.Status(fiber.StatusCreated).JSON(UserFromSQL(createdUser))
	})
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	MobileNumber int32  `json:"mobile_number"`
	UserType     string `json:"user_type"`
	BusinessID   string `json:"business_id,omitempty"`
	OutletID     string `json:"outlet_id,omitempty"`
	Password     string `json:"password"`
}

func UserFromSQL(u db.User) User {
	return User{
		ID:           u.ID,
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

func userTypeToString(i any) string {
	return fmt.Sprintf("%v", i)
}
