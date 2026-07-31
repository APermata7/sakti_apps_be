package domain

import "time"

type TelegramConnectRequest struct {
    VerificationCode string `json:"verification_code"`
}

type TelegramStatusResponse struct {
    IsConnected      bool   `json:"is_connected"`
    TelegramChatID   string `json:"telegram_chat_id,omitempty"`
    TelegramUsername string `json:"telegram_username,omitempty"`
}

type TelegramVerification struct {
    ID         string    `json:"id"`
    Code       string    `json:"code"`
    ChatID     string    `json:"chat_id"`
    Username   string    `json:"username"`
    KaryawanID *string   `json:"karyawan_id"`
    CreatedAt  time.Time `json:"created_at"`
    ExpiredAt  time.Time `json:"expired_at"`
    IsUsed     bool      `json:"is_used"`
    UsedAt     *time.Time `json:"used_at"`
}