package handler

import (
	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/usecase"
)

type TelegramHandler struct {
	TelegramUsecase *usecase.TelegramUsecase
}

func NewTelegramHandler(telegramUsecase *usecase.TelegramUsecase) *TelegramHandler {
	return &TelegramHandler{TelegramUsecase: telegramUsecase}
}

func (h *TelegramHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	status, err := h.TelegramUsecase.GetTelegramStatus(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    status,
	})
}

func (h *TelegramHandler) Connect(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req domain.TelegramConnectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.VerificationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kode verifikasi wajib diisi",
		})
	}

	if err := h.TelegramUsecase.ConnectTelegram(c.Context(), userID, req.VerificationCode); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Akun Telegram berhasil dihubungkan",
	})
}

func (h *TelegramHandler) Disconnect(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	if err := h.TelegramUsecase.DisconnectTelegram(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Koneksi Telegram berhasil diputuskan",
	})
}