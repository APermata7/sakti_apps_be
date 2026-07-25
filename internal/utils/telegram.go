package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TelegramBot struct {
	Token string
	DB    *pgxpool.Pool
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

func NewTelegramBot(db *pgxpool.Pool) *TelegramBot {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Println("TELEGRAM_BOT_TOKEN not found in environment")
		return nil
	}
	return &TelegramBot{
		Token: token,
		DB:    db,
	}
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

func (t *TelegramBot) SaveNotification(karyawanID, judul, pesan, referensiID, referensiTipe string) error {
	if t == nil || t.DB == nil {
		return nil
	}

	query := `
		INSERT INTO notifikasi (karyawan_id, jenis, channel, judul, pesan, referensi_id, referensi_tipe, dibuat_pada)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`

	_, err := t.DB.Exec(context.Background(), query,
		karyawanID,
		"pengajuan_cuti",
		"telegram",
		judul,
		pesan,
		referensiID,
		referensiTipe,
	)
	if err != nil {
		log.Printf("Failed to save notification to database: %v", err)
		return err
	}

	log.Printf("Telegram notification saved to database for karyawan: %s", karyawanID)
	return nil
}

func (t *TelegramBot) SendCreateLeaveNotification(chatID, karyawanID, karyawanNama, totalHari, alasan, leaveID string) error {
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

	err := t.SendMessage(chatID, text)
	if err != nil {
		return err
	}

	return t.SaveNotification(
		karyawanID,
		"Pengajuan Cuti Baru",
		karyawanNama+" mengajukan cuti "+totalHari+" hari",
		leaveID,
		"pengajuan_cuti",
	)
}

func (t *TelegramBot) SendCancelLeaveNotification(chatID, karyawanID, karyawanNama, leaveID string) error {
	sekarang := time.Now()
	tanggal := sekarang.Format("02 January 2006")

	text := fmt.Sprintf(
		"<b>📢 Pembatalan Pengajuan Cuti</b>\n\n"+
			"Yth. Atasan,\n\n"+
			"Pengajuan cuti berikut telah dibatalkan oleh karyawan:\n\n"+
			"👤 Karyawan : <b>%s</b>\n"+
			"📅 Tanggal Pembatalan : %s\n\n"+
			"📝 Alasan : Dibatalkan oleh karyawan\n\n"+
			"Tidak diperlukan proses persetujuan lebih lanjut.\n\n"+
			"Pesan ini dikirim secara otomatis oleh Sistem SAKTI.",
		karyawanNama, tanggal,
	)

	err := t.SendMessage(chatID, text)
	if err != nil {
		return err
	}

	return t.SaveNotification(
		karyawanID,
		"Pengajuan Cuti Dibatalkan",
		"Pengajuan cuti "+karyawanNama+" telah dibatalkan",
		leaveID,
		"pengajuan_cuti",
	)
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