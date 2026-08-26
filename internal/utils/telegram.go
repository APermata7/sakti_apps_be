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
	Ok          bool `json:"ok"`
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

	if chatID == "" || text == "" {
		return nil
	}

	if t.DB != nil {
		var status string
		query := `SELECT status_karyawan FROM karyawan WHERE telegram_chat_id = $1`
		err := t.DB.QueryRow(context.Background(), query, chatID).Scan(&status)
		if err != nil {
			return nil
		}
		if status != "aktif" {
			return nil
		}
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

func (t *TelegramBot) formatTanggalIndonesia(tgl time.Time) string {
	bulan := map[string]string{
		"January": "Januari", "February": "Februari", "March": "Maret",
		"April": "April", "May": "Mei", "June": "Juni",
		"July": "Juli", "August": "Agustus", "September": "September",
		"October": "Oktober", "November": "November", "December": "Desember",
	}
	bulanInggris := tgl.Format("January")
	bulanIndo := bulan[bulanInggris]
	return tgl.Format("02 " + bulanIndo + " 2006")
}

func (t *TelegramBot) formatTanggalRange(start, end time.Time) string {
	bulan := map[string]string{
		"January": "Januari", "February": "Februari", "March": "Maret",
		"April": "April", "May": "Mei", "June": "Juni",
		"July": "Juli", "August": "Agustus", "September": "September",
		"October": "Oktober", "November": "November", "December": "Desember",
	}
	bulanStart := bulan[start.Format("January")]
	bulanEnd := bulan[end.Format("January")]
	startStr := start.Format("02 " + bulanStart + " 2006")
	endStr := end.Format("02 " + bulanEnd + " 2006")
	return startStr + " - " + endStr
}

func (t *TelegramBot) SendCreateLeaveNotification(chatID, karyawanID, karyawanNama, totalHari, alasan, leaveID, tanggalMulai, tanggalSelesai string) error {
	tanggal := t.formatTanggalRange(
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
	)
	if tanggalMulai != "" && tanggalSelesai != "" {
		start, _ := time.Parse("2006-01-02", tanggalMulai)
		end, _ := time.Parse("2006-01-02", tanggalSelesai)
		tanggal = t.formatTanggalRange(start, end)
	}

	text := fmt.Sprintf(
		"<b>📋 Pengajuan Cuti Baru</b>\n\n"+
			"Yth. Atasan,\n\n"+
			"Seorang karyawan telah mengajukan cuti dan memerlukan persetujuan Anda. Berikut detail pengajuannya:\n\n"+
			"👤 Karyawan : <b>%s</b>\n"+
			"📅 Tanggal Cuti : <b>%s</b>\n"+
			"📅 Lama Cuti : <b>%s hari</b>\n"+
			"📝 Alasan : %s\n\n"+
			"Silakan lakukan proses persetujuan melalui aplikasi <b>SAKTI</b>.\n\n"+
			"Terima kasih.",
		karyawanNama, tanggal, totalHari, alasan,
	)

	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendFinalizationHRDNotification(chatID, karyawanID, karyawanNama, subTipe, tanggalCuti, leaveID string) error {
	text := fmt.Sprintf(
		"<b>📋 Pengajuan Cuti Baru</b>\n\n"+
			"Yth. HRD,\n\n"+
			"Seorang karyawan telah mengajukan cuti. Berikut detail pengajuannya:\n\n"+
			"👤 Karyawan : <b>%s</b>\n"+
			"📋 Jenis Cuti : <b>%s</b>\n"+
			"📅 Tanggal Cuti : <b>%s</b>\n\n"+
			"Silakan lakukan proses finalisasi melalui aplikasi <b>SAKTI</b>.\n\n"+
			"Terima kasih.",
		karyawanNama, subTipe, tanggalCuti,
	)

	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendCreateDispensasiHRDNotification(chatID, karyawanID, karyawanNama, totalHari, alasan, leaveID, tanggalMulai, tanggalSelesai string) error {
	tanggal := t.formatTanggalRange(
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
	)
	if tanggalMulai != "" && tanggalSelesai != "" {
		start, _ := time.Parse("2006-01-02", tanggalMulai)
		end, _ := time.Parse("2006-01-02", tanggalSelesai)
		tanggal = t.formatTanggalRange(start, end)
	}

	text := fmt.Sprintf(
		"<b>📋 Pengajuan Dispensasi Baru</b>\n\n"+
			"Yth. HRD,\n\n"+
			"Seorang karyawan telah mengajukan dispensasi. Berikut detail pengajuannya:\n\n"+
			"👤 Karyawan : <b>%s</b>\n"+
			"📅 Tanggal Dispensasi : <b>%s</b>\n"+
			"📅 Lama Dispensasi : <b>%s hari</b>\n"+
			"📝 Alasan : %s\n\n"+
			"Pengajuan ini telah langsung difinalisasi.\n\n"+
			"Terima kasih.",
		karyawanNama, tanggal, totalHari, alasan,
	)

	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendCreateDispensasiAtasanNotification(chatID, karyawanID, karyawanNama, totalHari, alasan, leaveID, tanggalMulai, tanggalSelesai string) error {
	tanggal := t.formatTanggalRange(
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
	)
	if tanggalMulai != "" && tanggalSelesai != "" {
		start, _ := time.Parse("2006-01-02", tanggalMulai)
		end, _ := time.Parse("2006-01-02", tanggalSelesai)
		tanggal = t.formatTanggalRange(start, end)
	}

	text := fmt.Sprintf(
		"<b>📋 Pengajuan Dispensasi Baru</b>\n\n"+
			"Yth. Atasan,\n\n"+
			"Seorang karyawan telah mengajukan dispensasi. Berikut detail pengajuannya:\n\n"+
			"👤 Karyawan : <b>%s</b>\n"+
			"📅 Tanggal Dispensasi : <b>%s</b>\n"+
			"📅 Lama Dispensasi : <b>%s hari</b>\n"+
			"📝 Alasan : %s\n\n"+
			"Pengajuan ini telah langsung disetujui.\n\n"+
			"Terima kasih.",
		karyawanNama, tanggal, totalHari, alasan,
	)

	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendCancelLeaveNotification(chatID, karyawanID, karyawanNama, leaveID, alasanBatal, tanggalMulai, tanggalSelesai string) error {
	tanggal := t.formatTanggalRange(
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
	)
	if tanggalMulai != "" && tanggalSelesai != "" {
		start, _ := time.Parse("2006-01-02", tanggalMulai)
		end, _ := time.Parse("2006-01-02", tanggalSelesai)
		tanggal = t.formatTanggalRange(start, end)
	}

	if alasanBatal == "" {
		alasanBatal = "Dibatalkan oleh karyawan"
	}

	text := fmt.Sprintf(
		"<b>📢 Pembatalan Pengajuan Cuti</b>\n\n"+
			"Yth. Atasan,\n\n"+
			"Pengajuan cuti berikut telah dibatalkan oleh karyawan:\n\n"+
			"👤 Karyawan : <b>%s</b>\n"+
			"📅 Tanggal Cuti : <b>%s</b>\n\n"+
			"📝 Alasan : %s\n\n"+
			"Tidak diperlukan proses persetujuan lebih lanjut.\n\n"+
			"Pesan ini dikirim secara otomatis oleh Sistem SAKTI.",
		karyawanNama, tanggal, alasanBatal,
	)

	return t.SendMessage(chatID, text)
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
			"Tanggal Cuti: %s\n\n"+
			"Silakan buka aplikasi SAKTI untuk detail lebih lanjut.",
		karyawanNama, status, tanggal,
	)
	return t.SendMessage(chatID, text)
}

func (t *TelegramBot) SendApprovalNotification(chatID, karyawanNama, totalHari, alasan, tanggalMulai, tanggalSelesai string) error {
	tanggal := t.formatTanggalRange(
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
	)
	if tanggalMulai != "" && tanggalSelesai != "" {
		start, _ := time.Parse("2006-01-02", tanggalMulai)
		end, _ := time.Parse("2006-01-02", tanggalSelesai)
		tanggal = t.formatTanggalRange(start, end)
	}

	text := fmt.Sprintf(
		"<b>Pengajuan Cuti Baru</b>\n\n"+
			"Karyawan: <b>%s</b>\n"+
			"Tanggal Cuti: <b>%s</b>\n"+
			"Lama Cuti: <b>%s hari</b>\n"+
			"Alasan: %s\n\n"+
			"Segera lakukan approval di aplikasi SAKTI.",
		karyawanNama, tanggal, totalHari, alasan,
	)
	return t.SendMessage(chatID, text)
}