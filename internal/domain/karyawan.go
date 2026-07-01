package domain

import "time"

type Karyawan struct {
    ID               string     `json:"id"`
    NamaLengkap      string     `json:"nama_lengkap"`
    Email            string     `json:"email"`
    NomorTelepon     *string    `json:"nomor_telepon"`    
    FotoURL          *string    `json:"foto_url"`         
    Peran            string     `json:"peran"`
    LevelJabatan     string     `json:"level_jabatan"`
    AtasanLangsungID *string    `json:"atasan_langsung_id"`
    Divisi           *string    `json:"divisi"`         
    Unit             *string    `json:"unit"`          
    StatusKaryawan   string     `json:"status_karyawan"`
    DibuatPada       time.Time  `json:"dibuat_pada"`
    DiperbaruiPada   time.Time  `json:"diperbarui_pada"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	User        Karyawan `json:"user"`
}

type CreateKaryawanRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	NamaLengkap      string `json:"nama_lengkap"`
	NomorTelepon     string `json:"nomor_telepon"`
	FotoURL          string `json:"foto_url"`
	Peran            string `json:"peran"`
	LevelJabatan     string `json:"level_jabatan"`
	AtasanLangsungID string `json:"atasan_langsung_id"`
	Divisi           string `json:"divisi"`
	Unit             string `json:"unit"`
}

type UpdateKaryawanRequest struct {
	NamaLengkap      string `json:"nama_lengkap"`
	NomorTelepon     string `json:"nomor_telepon"`
	FotoURL          string `json:"foto_url"`
	Peran            string `json:"peran"`
	LevelJabatan     string `json:"level_jabatan"`
	AtasanLangsungID string `json:"atasan_langsung_id"`
	Divisi           string `json:"divisi"`
	Unit             string `json:"unit"`
	StatusKaryawan   string `json:"status_karyawan"`
}