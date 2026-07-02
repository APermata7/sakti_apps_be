package domain

import "time"

type Presensi struct {
	ID              string    `json:"id"`
	KaryawanID      string    `json:"karyawan_id"`
	Tanggal         string    `json:"tanggal"`
	JamMasuk        string    `json:"jam_masuk"`
	JamKeluar       string    `json:"jam_keluar"`
	Status          string    `json:"status"`
	LintangMasuk    float64   `json:"lintang_masuk"`
	BujurMasuk      float64   `json:"bujur_masuk"`
	LintangKeluar   float64   `json:"lintang_keluar"`
	BujurKeluar     float64   `json:"bujur_keluar"`
	ValidasiWajah   bool      `json:"validasi_wajah"`
	FaceSimilarity  float64   `json:"face_similarity"`
	URLFoto         string    `json:"url_foto"`
	AlasanTerlambat string    `json:"alasan_terlambat"`
	Lembur          bool      `json:"lembur"`
	JamLembur       float64   `json:"jam_lembur"`
	DistanceMeter   float64   `json:"distance_meter"`
	IsOutsideRadius bool      `json:"is_outside_radius"`
	LocationStatus  string    `json:"location_status"`
	DibuatPada      time.Time `json:"dibuat_pada"`
	DiperbaruiPada  time.Time `json:"diperbarui_pada"`
}

type CheckInRequest struct {
	SelfieURL       string  `json:"selfie_url"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	AlasanTerlambat string  `json:"alasan_terlambat"`
}

type CheckOutRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CheckInResponse struct {
	ID              string  `json:"id"`
	KaryawanID      string  `json:"karyawan_id"`
	Tanggal         string  `json:"tanggal"`
	JamMasuk        string  `json:"jam_masuk"`
	Status          string  `json:"status"`
	ValidasiWajah   bool    `json:"validasi_wajah"`
	FaceSimilarity  float64 `json:"face_similarity"`
	URLFoto         string  `json:"url_foto"`
	DistanceMeter   float64 `json:"distance_meter"`
	IsOutsideRadius bool    `json:"is_outside_radius"`
	LocationStatus  string  `json:"location_status"`
	OfficeLatitude  float64 `json:"office_latitude"`
	OfficeLongitude float64 `json:"office_longitude"`
	OfficeRadius    int     `json:"office_radius"`
}

type CheckOutResponse struct {
	ID              string  `json:"id"`
	KaryawanID      string  `json:"karyawan_id"`
	Tanggal         string  `json:"tanggal"`
	JamMasuk        string  `json:"jam_masuk"`
	JamKeluar       string  `json:"jam_keluar"`
	Lembur          bool    `json:"lembur"`
	JamLembur       float64 `json:"jam_lembur"`
	DistanceMeter   float64 `json:"distance_meter"`
	IsOutsideRadius bool    `json:"is_outside_radius"`
	LocationStatus  string  `json:"location_status"`
}

type TodayResponse struct {
	HasCheckedIn   bool   `json:"has_checked_in"`
	HasCheckedOut  bool   `json:"has_checked_out"`
	CheckInTime    string `json:"check_in_time"`
	CheckInStatus  string `json:"check_in_status"`
	Tanggal        string `json:"tanggal"`
	KaryawanID     string `json:"karyawan_id"`
}

type FaceVerificationRequest struct {
	SelfieURL   string `json:"selfie_url"`
	KaryawanID  string `json:"karyawan_id"`
}