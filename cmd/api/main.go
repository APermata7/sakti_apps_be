package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env file not found, using system env")
    }

    app := fiber.New()

    app.Use(logger.New())
    app.Use(cors.New())

    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status":  "ok",
            "service": "sakti-api",
        })
    })

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 Server running on port %s", port)
    if err := app.Listen(":" + port); err != nil {
        log.Fatal(err)
    }
}