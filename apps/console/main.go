package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
)

//go:embed index.html
var consoleHTML []byte

type ConsoleView struct {
	Email string
}

func main() {
	app := fiber.New()

	tmpl, err := template.New("console").Parse(string(consoleHTML))
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		session := c.Cookies("user_session")

		if session == "" {
			return c.Redirect("http://localhost:3001/login?error=Please%20login%20to%20access%20console")
		}

		var buf bytes.Buffer

		view := ConsoleView{
			Email: session,
		}

		if err := tmpl.Execute(&buf, view); err != nil {
			return err
		}

		c.Set("Content-Type", "text/html; charset=utf-8")

		return c.Send(buf.Bytes())
	})

	log.Println("Console running on http://localhost:3002")

	log.Fatal(app.Listen(":3002"))
}