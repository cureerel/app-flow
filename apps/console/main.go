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

	log.Fatal(app.Listen(":3002"))
}