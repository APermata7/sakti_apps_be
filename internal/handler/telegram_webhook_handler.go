package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sakti_apps_be/internal/usecase"
)

type TelegramWebhookHandler struct {
	TelegramUsecase *usecase.TelegramUsecase
}

func NewTelegramWebhookHandler(telegramUsecase *usecase.TelegramUsecase) *TelegramWebhookHandler {
	return &TelegramWebhookHandler{TelegramUsecase: telegramUsecase}
}

type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID int `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

func (h *TelegramWebhookHandler) Webhook(c *fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		log.Println("[Webhook] Empty body")
		return c.Status(400).SendString("Bad Request")
	}

	var update TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		log.Printf("[Webhook] Parse error: %v", err)
		return c.Status(400).SendString("Bad Request")
	}

	log.Printf("[Webhook] Received: text=%s, chat_id=%d", update.Message.Text, update.Message.Chat.ID)

	if update.Message.Text == "/start" {
		chatID := strconv.Itoa(update.Message.Chat.ID)
		username := update.Message.From.Username

		code, err := h.TelegramUsecase.GenerateVerificationCode(c.Context(), chatID, username)
		if err != nil {
			log.Printf("[Webhook] Generate code error: %v", err)
			return c.Status(500).SendString("Internal Server Error")
		}

		text := "👋 Selamat datang di SAKTI Bot!\n\n" +
			"Kode verifikasi Anda: <b>" + code + "</b>\n\n" +
			"Masukkan kode ini di aplikasi SAKTI untuk menghubungkan akun Telegram.\n\n" +
			"⏰ Kode berlaku selama 5 menit.\n\n" +
			"Jika tidak digunakan, abaikan pesan ini."

		token := os.Getenv("TELEGRAM_BOT_TOKEN")
		if token == "" {
			log.Println("[Webhook] TELEGRAM_BOT_TOKEN not found")
			return c.Status(500).SendString("Internal Server Error")
		}

		sendURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
		reqBody := map[string]interface{}{
			"chat_id":    chatID,
			"text":       text,
			"parse_mode": "HTML",
		}
		jsonBody, _ := json.Marshal(reqBody)

		resp, err := http.Post(sendURL, "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Printf("[Webhook] Send message error: %v", err)
			return c.Status(500).SendString("Internal Server Error")
		}
		defer resp.Body.Close()

		log.Printf("[Webhook] Code sent to chat_id: %s, code: %s", chatID, code)
	}

	return c.Status(200).SendString("OK")
}