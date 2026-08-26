package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/usecase"
)

type TelegramHandler struct {
	TelegramUsecase *usecase.TelegramUsecase
}

func NewTelegramHandler(telegramUsecase *usecase.TelegramUsecase) *TelegramHandler {
	return &TelegramHandler{TelegramUsecase: telegramUsecase}
}

func (h *TelegramHandler) GetStatus(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Println("[TelegramStatus] user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	status, err := h.TelegramUsecase.GetTelegramStatus(c.Context(), userID)
	if err != nil {
		log.Printf("[TelegramStatus] error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	connected, _ := status["connected"].(bool)
	chatID, _ := status["chat_id"].(string)

	log.Printf("[TelegramStatus] userID=%s, connected=%v, chatID=%s", userID, connected, chatID)

	response := fiber.Map{
		"success": true,
		"data":    status,
	}

	log.Printf("[TelegramStatus] full response: %+v", response)

	return c.JSON(response)
}

func (h *TelegramHandler) Connect(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Println("[ConnectTelegram] user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var req struct {
		VerificationCode string `json:"verification_code"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ConnectTelegram] body parse error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	log.Printf("[ConnectTelegram] userID=%s, code=%s", userID, req.VerificationCode)

	if req.VerificationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kode verifikasi wajib diisi",
		})
	}

	if err := h.TelegramUsecase.ConnectTelegram(c.Context(), userID, req.VerificationCode); err != nil {
		log.Printf("[ConnectTelegram] error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	log.Printf("[ConnectTelegram] success for userID=%s", userID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Telegram berhasil terhubung",
	})
}

func (h *TelegramHandler) Disconnect(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Println("[DisconnectTelegram] user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	log.Printf("[DisconnectTelegram] userID=%s", userID)

	if err := h.TelegramUsecase.ClearChatID(c.Context(), userID); err != nil {
		log.Printf("[DisconnectTelegram] error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	log.Printf("[DisconnectTelegram] success for userID=%s", userID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Koneksi Telegram berhasil diputuskan",
	})
}