package middlewares

import (
	"module_6/internal/utils"

	"github.com/gofiber/fiber/v3"
)

func JWTProtected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		tokenStr := authHeader[len("Bearer "):]
		token, err := utils.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		return c.Next()
	}
}
