package handler

import (
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if config == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Konfigurasi tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    config,
	})
}

func (h *KonfigurasiHandler) GetWorkConfig(c *fiber.Ctx) error {
	config, err := h.KonfigurasiUsecase.GetWorkConfig(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if config == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Konfigurasi tidak ditemukan",
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

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca file",
		})
	}
	defer src.Close()

	config, err := h.KonfigurasiUsecase.UploadLogo(c.Context(), userID, src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logo berhasil diupload",
		"data":    config,
	})
}