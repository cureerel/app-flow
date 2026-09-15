package main

import (
    "bytes"
    _ "embed"
    "log"
    "net/url"
    "text/template"

    "github.com/gofiber/fiber/v2"
)

//go:embed index.html
var indexHTML []byte

type AuthView struct {
    Content string
}

func main() {
    app := fiber.New()
    tmpl, err := template.New("index").Parse(string(indexHTML))
    if err != nil {
        log.Fatal(err)
    }

    renderPage := func(c *fiber.Ctx, content string, errorMsg string) error {
        if errorMsg != "" {
            // SVG Icon + Error Message
            errorMsg = `
            <div class="error-msg">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
                <span>` + errorMsg + `</span>
            </div>`
            content = errorMsg + content
        }

        var buf bytes.Buffer
        if err := tmpl.Execute(&buf, AuthView{Content: content}); err != nil {
            return err
        }
        c.Set("Content-Type", "text/html")
        return c.Send(buf.Bytes())
    }

    app.Get("/", func(c *fiber.Ctx) error {
        return c.Redirect("/login")
    })

    app.Get("/login", func(c *fiber.Ctx) error {
        errorMsg := c.Query("error", "")
        if decoded, err := url.QueryUnescape(errorMsg); err == nil {
            errorMsg = decoded
        }

        loginHTML := `
            <h1>Welcome back</h1>
            <p style="color:#6B7280; margin-top:-16px; margin-bottom:24px;">Enter your details to access your account.</p>
            
            <form action="http://localhost:4000/auth/login" method="POST">
                <input type="email" name="email" placeholder="Email address" required autocomplete="email">
                <input type="password" name="password" placeholder="Password" required autocomplete="current-password">
                
                <button type="submit" class="btn btn-primary">Log in</button>
            </form>

            <div class="htmx-indicator">Verifying credentials...</div>

            <div style="text-align: center; margin-top: 24px;">
                <span style="font-size: 0.9rem; color: #6B7280;">Don't have an account?</span><br>
                <a href="/signup" class="btn btn-outline" style="margin-top: 8px;">Sign up</a>
            </div>
        `
        return renderPage(c, loginHTML, errorMsg)
    })

    app.Get("/signup", func(c *fiber.Ctx) error {
        errorMsg := c.Query("error", "")
        if decoded, err := url.QueryUnescape(errorMsg); err == nil {
            errorMsg = decoded
        }

        signupHTML := `
            <h1>Create account</h1>
            <p style="color:#6B7280; margin-top:-16px; margin-bottom:24px;">Get started with MyApp today.</p>
            
            <form action="http://localhost:4000/auth/signup" method="POST">
                <input type="text" name="name" placeholder="Full Name" required autocomplete="name">
                <input type="email" name="email" placeholder="Email address" required autocomplete="email">
                <input type="password" name="password" placeholder="Create Password" required autocomplete="new-password">
                
                <button type="submit" class="btn btn-primary">Create account</button>
            </form>

            <div class="htmx-indicator">Creating account...</div>

            <div style="text-align: center; margin-top: 24px;">
                <span style="font-size: 0.9rem; color: #6B7280;">Already have an account?</span><br>
                <a href="/login" class="btn btn-outline" style="margin-top: 8px;">Log in</a>
            </div>
        `
        return renderPage(c, signupHTML, errorMsg)
    })

    log.Fatal(app.Listen(":3001"))
}