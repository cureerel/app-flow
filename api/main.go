package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Post("/auth/login", func(c *fiber.Ctx) error {
		var req AuthRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"ok":    false,
				"error": "invalid request",
			})
		}

		return c.JSON(fiber.Map{
			"ok":    true,
			"email": req.Email,
		})
	})

	app.Post("/auth/signup", func(c *fiber.Ctx) error {
		var req AuthRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"ok":    false,
				"error": "invalid request",
			})
		}

		return c.JSON(fiber.Map{
			"ok":    true,
			"email": req.Email,
		})
	})

	log.Fatal(app.Listen(":4000"))
}