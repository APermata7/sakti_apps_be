package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/usecase"
)

type KonfigurasiHandler struct {
	KonfigurasiUsecase *usecase.KonfigurasiUsecase
}

func NewKonfigurasiHandler(konfigurasiUsecase *usecase.KonfigurasiUsecase) *KonfigurasiHandler {
	return &KonfigurasiHandler{KonfigurasiUsecase: konfigurasiUsecase}
}

func (h *KonfigurasiHandler) GetConfig(c *fiber.Ctx) error {
	config, err := h.KonfigurasiUsecase.GetConfig(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    config,
	})
}

func (h *KonfigurasiHandler) UpdateConfig(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req domain.UpdateKonfigurasiRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	config, err := h.KonfigurasiUsecase.UpdateConfig(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Konfigurasi berhasil diupdate",
		"data":    config,
	})
}

func (h *KonfigurasiHandler) UploadLogo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	file, err := c.FormFile("logo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File logo wajib diupload",
		})
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format file harus JPG, PNG, atau WEBP",
		})
	}

	if file.Size > 2*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ukuran file maksimal 2MB",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuka file",
		})
	}
	defer src.Close()

	filename := "logo_kantor_" + time.Now().Format("20060102150405")

	config, err := h.KonfigurasiUsecase.UploadLogo(c.Context(), userID, src, filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logo berhasil diupload",
		"data": fiber.Map{
			"logo_url": config.LogoKantor,
			"config":   config,
		},
	})
}