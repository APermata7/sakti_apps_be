package domain

import "time"

type FCMToken struct {
	ID          string    `json:"id"`
	KaryawanID  string    `json:"karyawan_id"`
	FCMToken    string    `json:"fcm_token"`
	IsActive    bool      `json:"is_active"`
	DibuatPada  time.Time `json:"dibuat_pada"`
	DiperbaruiPada time.Time `json:"diperbarui_pada"`
}

type Notifikasi struct {
	ID            string     `json:"id"`
	KaryawanID    string     `json:"karyawan_id"`
	Jenis         string     `json:"jenis"`
	Channel       string     `json:"channel"`
	Judul         string     `json:"judul"`
	Pesan         string     `json:"pesan"`
	Dibaca        bool       `json:"dibaca"`
	DibacaPada    *time.Time `json:"dibaca_pada"`
	ReferensiID   string     `json:"referensi_id"`
	ReferensiTipe string     `json:"referensi_tipe"`
	DibuatPada    time.Time  `json:"dibuat_pada"`
}

type KirimNotifikasiRequest struct {
	KaryawanID    string `json:"karyawan_id"`
	Jenis         string `json:"jenis"`
	Judul         string `json:"judul"`
	Pesan         string `json:"pesan"`
	ReferensiID   string `json:"referensi_id"`
	ReferensiTipe string `json:"referensi_tipe"`
}

type NotifikasiResponse struct {
	Items []Notifikasi   `json:"items"`
	Meta  MetaPagination `json:"meta"`
}