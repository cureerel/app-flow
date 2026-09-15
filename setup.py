from pathlib import Path
import subprocess
import sys


ROOT = Path.cwd()

APPS = {
    "landing": 3000,
    "auth": 3001,
    "console": 3002,
    "billing": 3003,
}

FILES = {
    ".gitignore": """\
.DS_Store
.env
*.sqlite
*.sqlite3
bin/
""",

    "README.md": """\
# Multi App Test Harness

## Apps

- landing: http://localhost:3000
- auth: http://localhost:3001
- console: http://localhost:3002
- billing: http://localhost:3003
- api: http://localhost:4000

## Stack

- Go
- Fiber
- HTMX
- SQLite
""",

    "db/.gitkeep": "",
    "tests/.gitkeep": "",

    "api/main.go": """\
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
""",
}


def write_file(path: str, content: str):
    target = ROOT / path
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(content)


def create_app(name: str, port: int):
    app_dir = ROOT / "apps" / name
    app_dir.mkdir(parents=True, exist_ok=True)

    main = f'''package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
)

func main() {{
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {{
        return c.SendString("{name} app")
    }})

    log.Fatal(app.Listen(":{port}"))
}}
'''

    html = f'''<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{name}</title>
    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
</head>
<body>
    <main>
        <h1>{name}</h1>
    </main>
</body>
</html>
'''

    write_file(f"apps/{name}/main.go", main)
    write_file(f"apps/{name}/index.html", html)


def run(command):
    print("$", " ".join(command))
    subprocess.run(command, cwd=ROOT, check=True)


def main():
    print(f"Setting up: {ROOT}")

    if any((ROOT / x).exists() for x in ["go.mod", "apps", "api"]):
        print("Warning: this folder already contains project files.")

        answer = input("Continue? [y/N] ").strip().lower()

        if answer != "y":
            sys.exit(1)

    for directory in [
        "apps",
        "api",
        "db",
        "tests",
    ]:
        (ROOT / directory).mkdir(parents=True, exist_ok=True)

    for path, content in FILES.items():
        write_file(path, content)

    for name, port in APPS.items():
        create_app(name, port)

    run(["go", "mod", "init", "multi-test"])
    run(["go", "get", "github.com/gofiber/fiber/v2"])
    run(["go", "mod", "tidy"])

    print()
    print("Setup complete.")
    print()
    print("Apps:")
    for name, port in APPS.items():
        print(f"  {name:10} http://localhost:{port}")

    print("  api        http://localhost:4000")
    print()
    print("Run an app with:")
    print("  go run ./apps/landing")
    print("  go run ./apps/auth")
    print("  go run ./apps/console")
    print("  go run ./apps/billing")
    print("  go run ./api")


if __name__ == "__main__":
    main()