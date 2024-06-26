package middleware

import (
	"github.com/paimon_bank/internal/customErr"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func NewAuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return customErr.NewUnauthorizedError("Unauthorized")
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)

			userLoggedInId := claims["ID"].(string)
			c.Locals("userLoggedInId", userLoggedInId)

			return c.Next()
		},
	})
}
