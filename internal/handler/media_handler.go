package handler

import (
	"log"

	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open file",
		})
	}
	defer src.Close()

	url, err := utils.UploadFile(src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"url":     url,
	})
}

func UploadImage(c *fiber.Ctx) error {
	log.Println("UploadImage: start")

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	karyawanID := c.FormValue("karyawan_id")
	log.Printf("UploadImage: role=%s, userID=%s, karyawanID=%s", role, userID, karyawanID)

	if !(role == "admin" && karyawanID != "") {
		karyawanID = userID
	}
	log.Printf("UploadImage: final karyawanID=%s", karyawanID)

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("UploadImage: FormFile error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Image not found",
		})
	}
	log.Printf("UploadImage: file=%s, size=%d", file.Filename, file.Size)

	src, err := file.Open()
	if err != nil {
		log.Printf("UploadImage: Open error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open image",
		})
	}
	defer src.Close()

	log.Println("UploadImage: uploading to Cloudinary")
	url, err := utils.UploadImage(src, file.Filename)
	if err != nil {
		log.Printf("UploadImage: Cloudinary error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal upload: " + err.Error(),
		})
	}
	log.Printf("UploadImage: Cloudinary URL: %s", url)

	dbPool, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok || dbPool == nil {
		log.Printf("UploadImage: Database connection not found")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database connection not found",
		})
	}
	log.Println("UploadImage: database connection OK")

	karyawanRepo := repository.NewKaryawanRepo(dbPool)
	err = karyawanRepo.UpdateFotoURL(c.Context(), karyawanID, url)
	if err != nil {
		log.Printf("UploadImage: UpdateFotoURL error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal simpan ke database: " + err.Error(),
		})
	}
	log.Println("UploadImage: success")

	return c.JSON(fiber.Map{
		"success":     true,
		"url":         url,
		"karyawan_id": karyawanID,
	})
}