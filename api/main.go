package main

import (
    "database/sql"
    "log"
    "net/url"
    "time"

    "github.com/gofiber/fiber/v2"
    _ "modernc.org/sqlite"
)

type AuthRequest struct {
    Email    string `form:"email"`
    Password string `form:"password"`
    Name     string `form:"name"`
}

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite", "./db/app.db")
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }
    defer db.Close()

    createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME
    );`
    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatal("Failed to create table:", err)
    }

    app := fiber.New()


    app.Get("/health", func(c *fiber.Ctx) error {
	log.Println("GET /health")

	return c.JSON(fiber.Map{
		"ok": true,
	})
    })

    app.Post("/auth/signup", func(c *fiber.Ctx) error {
	log.Println("POST /auth/signup")

	var req AuthRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println("BodyParser error:", err)
		return redirectWithError(c, "/signup", "Invalid request")
	}

	log.Println("Signup:", req.Email)

	_, err := db.Exec(
		"INSERT INTO users (name, email, password, created_at) VALUES (?, ?, ?, ?)",
		req.Name,
		req.Email,
		req.Password,
		time.Now(),
	)

	if err != nil {
		log.Println("Signup DB error:", err)
		return redirectWithError(c, "/signup", "Email is already registered")
	}

	log.Println("Signup successful:", req.Email)

	setAuthCookie(c, req.Email)

	return c.Redirect("http://localhost:3002/")
    })

    // login

    app.Post("/auth/login", func(c *fiber.Ctx) error {
	log.Println("POST /auth/login")

	var req AuthRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println("BodyParser error:", err)
		return redirectWithError(c, "/login", "Invalid request")
	}

	log.Println("Login:", req.Email)

	var id int
	var email string

	err := db.QueryRow(
		"SELECT id, email FROM users WHERE email = ? AND password = ?",
		req.Email,
		req.Password,
	).Scan(&id, &email)

	if err != nil {
		log.Println("Login error:", err)
		return redirectWithError(c, "/login", "Invalid email or password")
	}

	log.Println("Login successful:", email)

	setAuthCookie(c, email)

	return c.Redirect("http://localhost:3002/")
    })

    // logout

    app.Post("/auth/logout", func(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "user_session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.Redirect("http://localhost:3001/login")
    })

    log.Fatal(app.Listen(":4000"))
}

func setAuthCookie(c *fiber.Ctx, email string) {
    c.Cookie(&fiber.Cookie{
        Name:     "user_session",
        Value:    email,
        Expires:  time.Now().Add(24 * time.Hour),
        HTTPOnly: true,
        SameSite: "Lax",
    })
}

func redirectWithError(c *fiber.Ctx, path string, msg string) error {
    return c.Redirect("http://localhost:3001" + path + "?error=" + url.QueryEscape(msg))
}