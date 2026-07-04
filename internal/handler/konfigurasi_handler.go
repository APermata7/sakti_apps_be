package handler

import (
	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/usecase"
	"sakti_apps_be/internal/utils"
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
		return c.Status(500).JSON(fiber.Map{
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
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	config, err := h.KonfigurasiUsecase.UpdateConfig(c.Context(), userID, req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    config,
	})
}

func (h *KonfigurasiHandler) UploadLogo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	file, err := c.FormFile("logo")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "File logo tidak ditemukan",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuka file",
		})
	}
	defer src.Close()

	url, err := utils.UploadImage(src, "logo_kopegtel")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal upload ke Cloudinary: " + err.Error(),
		})
	}

	config, err := h.KonfigurasiUsecase.UpdateLogo(c.Context(), userID, url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"logo_url": url,
			"config":   config,
		},
	})
}