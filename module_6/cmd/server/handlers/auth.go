package handlers

import (
	"module_6/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(c fiber.Ctx) error {
	u := new(User)
	if err := c.Bind().Body(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	return c.JSON(fiber.Map{"message": "user registered"})
}

func SignIn(c fiber.Ctx) error {
	u := new(User)
	if err := c.Bind().Body(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	token, err := utils.GenerateToken(u.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
	}
	utils.Logger.Println("Logged successfully:", u.Username)
	return c.JSON(fiber.Map{"token": token})
}
