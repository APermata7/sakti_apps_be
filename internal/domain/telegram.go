package domain

type TelegramConnectRequest struct {
	VerificationCode string `json:"verification_code"`
}

type TelegramStatusResponse struct {
	IsConnected      bool   `json:"is_connected"`
	TelegramChatID   string `json:"telegram_chat_id,omitempty"`
	TelegramUsername string `json:"telegram_username,omitempty"`
}

type TelegramVerificationCode struct {
	Code      string `json:"code"`
	ChatID    string `json:"chat_id"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
}