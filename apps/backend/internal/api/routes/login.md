```
package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/your-project/db"
	"github.com/your-project/utils"
)

type AuthHandler struct {
	Queries *db.Queries
	Pool    *pgxpool.Pool
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get user from database
	user, err := h.Queries.GetUserByUsernameOrEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check password
	if !utils.ComparePasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check if valid session exists
	session, err := h.Queries.GetLatestUserSession(c.Context(), user.ID)
	if err == nil && session.ExpireAt.Time.After(time.Now()) {
		// Valid session exists, return existing tokens
		return c.JSON(LoginResponse{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
		})
	}

	// Generate new tokens
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
	}

	// Start a transaction
	tx, err := h.Pool.Begin(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}
	defer tx.Rollback(c.Context())

	qtx := h.Queries.WithTx(tx)

	// Delete existing sessions for the user
	if err := qtx.DeleteUserSessions(c.Context(), user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete old sessions"})
	}

	// Create new session
	newSession, err := qtx.CreateUserSession(c.Context(), db.CreateUserSessionParams{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt: db.Timestamp{
			Time:  time.Now().Add(30 * time.Minute),
			Valid: true,
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
	}

	// Commit the transaction
	if err := tx.Commit(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.JSON(LoginResponse{
		AccessToken:  newSession.AccessToken,
		RefreshToken: newSession.RefreshToken,
	})
}
```


postgres query for sqlc

```
-- name: GetLatestUserSession :one
SELECT * FROM user_sessions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteUserSessions :exec
DELETE FROM user_sessions
WHERE user_id = $1;
```