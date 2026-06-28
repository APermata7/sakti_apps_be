package domain

import "time"

type FCMToken struct {
	ID         string    `json:"id"`
	KaryawanID string    `json:"karyawan_id"`
	FCMToken   string    `json:"fcm_token"`
	DeviceID   string    `json:"device_id"`
	DeviceType string    `json:"device_type"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NotificationRequest struct {
	KaryawanID string            `json:"karyawan_id"`
	Title      string            `json:"title"`
	Body       string            `json:"body"`
	Data       map[string]string `json:"data,omitempty"`
}

type NotificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}