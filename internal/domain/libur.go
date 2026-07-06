package domain

import "time"

type Libur struct {
	ID        string    `json:"id"`
	Tanggal   string    `json:"tanggal"`
	Nama      string    `json:"nama"`
	Jenis     string    `json:"jenis"`
	Aktif     bool      `json:"aktif"`
	Sumber    string    `json:"sumber"`
	DibuatPada time.Time `json:"dibuat_pada"`
	DiperbaruiPada time.Time `json:"diperbarui_pada"`
}

type CreateLiburRequest struct {
	Tanggal string `json:"tanggal"`
	Nama    string `json:"nama"`
	Jenis   string `json:"jenis"`
}

type UpdateLiburRequest struct {
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Aktif bool   `json:"aktif"`
}