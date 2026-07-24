package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type TelegramBot struct {
	Token string
}

type SendMessageRequest struct {
	ChatID              string `json:"chat_id"`
	Text                string `json:"text"`
	ParseMode           string `json:"parse_mode,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
}

type SendMessageResponse struct {
	Ok          bool   `json:"ok"`
	Result      struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID int `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func NewTelegramBot() *TelegramBot {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Println("TELEGRAM_BOT_TOKEN not found in environment")
		return nil
	}
	return &TelegramBot{Token: token}
}

func (t *TelegramBot) SendMessage(chatID, text string) error {
	if t == nil {
		return fmt.Errorf("telegram bot not initialized")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)

	reqBody := SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.Ok {
		return fmt.Errorf("telegram error: %s", response.Description)
	}

	log.Printf("Telegram message sent to %s", chatID)
	return nil
}

func (t *TelegramBot) SendNotification(chatID, title, message string) error {
	text := fmt.Sprintf("<b>%s</b>\n\n%s", title, message)
	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendLeaveNotification(chatID, karyawanNama, status, tanggal string) error {
	text := fmt.Sprintf(
		"<b>Notifikasi Cuti</b>\n\n"+
			"Karyawan: <b>%s</b>\n"+
			"Status: <b>%s</b>\n"+
			"Tanggal: %s\n\n"+
			"Silakan buka aplikasi SAKTI untuk detail lebih lanjut.",
		karyawanNama, status, tanggal,
	)
	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendApprovalNotification(chatID, karyawanNama, totalHari, alasan string) error {
	text := fmt.Sprintf(
		"<b>Pengajuan Cuti Baru</b>\n\n"+
			"Karyawan: <b>%s</b>\n"+
			"Lama Cuti: <b>%s hari</b>\n"+
			"Alasan: %s\n\n"+
			"Segera lakukan approval di aplikasi SAKTI.",
		karyawanNama, totalHari, alasan,
	)
	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendCreateLeaveNotification(chatID, karyawanNama, totalHari, alasan string) error {
	text := fmt.Sprintf(
		"<b>📋 Pengajuan Cuti Baru</b>\n\n"+
			"Yth. Atasan,\n\n"+
			"Seorang karyawan telah mengajukan cuti dan memerlukan persetujuan Anda. Berikut detail pengajuannya:\n\n"+
			"👤 Karyawan : <b>%s</b>\n"+
			"📅 Lama Cuti : <b>%s hari</b>\n"+
			"📝 Alasan : %s\n\n"+
			"Silakan lakukan proses persetujuan melalui aplikasi <b>SAKTI</b> sesuai kebijakan yang berlaku.\n\n"+
			"Terima kasih.",
		karyawanNama, totalHari, alasan,
	)
	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendCancelLeaveNotification(chatID, karyawanNama, alasan string) error {
	sekarang := time.Now()
	tanggal := sekarang.Format("02 January 2006")

	text := fmt.Sprintf(
		"<b>📢 Pembatalan Pengajuan Cuti</b>\n\n"+
			"Yth. Atasan,\n\n"+
			"Pengajuan cuti berikut telah dibatalkan oleh karyawan:\n\n"+
			"👤 Karyawan : <b>%s</b>\n"+
			"📅 Tanggal Pembatalan : %s\n\n"+
			"📝 Alasan : %s\n\n"+
			"Tidak diperlukan proses persetujuan lebih lanjut.\n\n"+
			"Pesan ini dikirim secara otomatis oleh Sistem SAKTI.",
		karyawanNama, tanggal, alasan,
	)
	return t.SendMessage(chatID, text)
}