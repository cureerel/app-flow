package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "ok": true,
        })
    })

    log.Fatal(app.Listen(":4000"))
}
