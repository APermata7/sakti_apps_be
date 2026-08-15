package handler

import (
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
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	var req struct {
		KaryawanID string `json:"karyawan_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		req.KaryawanID = ""
	}

	var karyawanID string
	if role == "admin" && req.KaryawanID != "" {
		karyawanID = req.KaryawanID
	} else {
		karyawanID = userID
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Image not found",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open image",
		})
	}
	defer src.Close()

	url, err := utils.UploadImage(src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	dbPool, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database connection not found",
		})
	}

	karyawanRepo := repository.NewKaryawanRepo(dbPool)
	err = karyawanRepo.UpdateFotoURL(c.Context(), karyawanID, url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal simpan ke database: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"url":         url,
		"karyawan_id": karyawanID,
	})
}