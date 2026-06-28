package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"sakti_apps_be/internal/handler"
	"sakti_apps_be/internal/utils"
	"sakti_apps_be/pkg/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	log.Println("Menghubungkan ke database...")
	dbConn, err := db.NewSupabaseDB()
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer dbConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbConn.Ping(ctx); err != nil {
		log.Fatalf("Database tidak merespon: %v", err)
	}
	log.Println("Database connected")

	if err := utils.InitCloudinary(); err != nil {
		log.Printf("Cloudinary init error: %v", err)
	} else {
		log.Println("Cloudinary ready")
	}

	if err := utils.InitFCM(); err != nil {
		log.Printf("FCM init error: %v", err)
	} else {
		log.Println("FCM ready")
	}

	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": os.Getenv("APP_NAME"),
			"db":      "connected",
		})
	})

	app.Post("/upload/file", handler.UploadFile)
	app.Post("/upload/image", handler.UploadImage)
	app.Post("/upload/ttd", handler.UploadTTD)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}