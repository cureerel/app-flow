package main

import (
	_ "embed"
	"log"

	"github.com/gofiber/fiber/v2"
)

//go:embed index.html
var indexHTML []byte

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.Send(indexHTML)
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.Send(indexHTML)
	})

	app.Get("/signup", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.Send(indexHTML)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return c.Redirect("http://localhost:3002/")
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		return c.Redirect("http://localhost:3002/")
	})

	log.Fatal(app.Listen(":3001"))
}