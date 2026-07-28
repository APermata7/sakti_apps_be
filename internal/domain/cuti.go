package domain

import "time"

type PengajuanCuti struct {
	ID                  string     `json:"id"`
	KaryawanID          string     `json:"karyawan_id"`
	SubTipe             string     `json:"sub_tipe"`
	TanggalMulai        time.Time  `json:"tanggal_mulai"`
	TanggalSelesai      time.Time  `json:"tanggal_selesai"`
	TotalHari           int        `json:"total_hari"`
	Alasan              string     `json:"alasan"`
	Status              string     `json:"status"`
	BackDate            bool       `json:"back_date"`
	MengurangiCuti      bool       `json:"mengurangi_cuti"`
	LangsungApprove     bool       `json:"langsung_approve"`
	LangsungFinal       bool       `json:"langsung_final"`
	JudulDokumen        string     `json:"judul_dokumen"`
	DisetujuiOleh       *string    `json:"disetujui_oleh"`
	TanggalDisetujui    *time.Time `json:"tanggal_disetujui"`
	DifinalisasiOleh    *string    `json:"difinalisasi_oleh"`
	TanggalDifinalisasi *time.Time `json:"tanggal_difinalisasi"`
	URLPDF              *string    `json:"url_pdf"`
	AlasanBatal         *string    `json:"alasan_batal"`
	TanggalDibatalkan   *time.Time `json:"tanggal_dibatalkan"`
	AlasanDitolak       *string    `json:"alasan_ditolak"`
	TanggalDitolak      *time.Time `json:"tanggal_ditolak"`
	DibuatPada          time.Time  `json:"dibuat_pada"`
	DiperbaruiPada      time.Time  `json:"diperbarui_pada"`
}

type SisaCuti struct {
	ID                string    `json:"id"`
	KaryawanID        string    `json:"karyawan_id"`
	Tahun             int       `json:"tahun"`
	JumlahCuti        int       `json:"jumlah_cuti"`
	TelahDilaksanakan int       `json:"telah_dilaksanakan"`
	AkanDilaksanakan  int       `json:"akan_dilaksanakan"`
	SisaCuti          int       `json:"sisa_cuti"`
	DibuatPada        time.Time `json:"dibuat_pada"`
	DiperbaruiPada    time.Time `json:"diperbarui_pada"`
}

type CreateCutiRequest struct {
	SubTipe         string `json:"sub_tipe"`
	TanggalMulai    string `json:"tanggal_mulai"`
	TanggalSelesai  string `json:"tanggal_selesai"`
	Alasan          string `json:"alasan"`
	BackDate        bool   `json:"back_date"`
	LangsungApprove bool   `json:"langsung_approve"`
}

type ApproveCutiRequest struct{}

type RejectCutiRequest struct {
	Alasan string `json:"alasan"`
}

type FinalizeCutiRequest struct {
	Catatan string `json:"catatan"`
}

type BalanceResponse struct {
	Tahun                        int    `json:"tahun"`
	JumlahCuti                   int    `json:"jumlah_cuti"`
	TelahDilaksanakan            int    `json:"telah_dilaksanakan"`
	AkanDilaksanakan             int    `json:"akan_dilaksanakan"`
	SisaCuti                     int    `json:"sisa_cuti"`
	SisaCutiTahunIni             int    `json:"sisa_cuti_tahun_ini"`
	SisaCutiTahunLalu            int    `json:"sisa_cuti_tahun_lalu"`
	TotalCutiTersedia            int    `json:"total_cuti_tersedia"`
	KuotaPengajuanTersedia       int    `json:"kuota_pengajuan_tersedia"`
	SisaCutiTahunLaluBerlakuSampai string `json:"sisa_cuti_tahun_lalu_berlaku_sampai,omitempty"`
}

type LeaveFilterRequest struct {
	Status    string `json:"status"`
	SubTipe   string `json:"sub_tipe"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
}

type LeaveWithKaryawanResponse struct {
	PengajuanCuti
	KaryawanNama   string `json:"karyawan_nama"`
	KaryawanDivisi string `json:"karyawan_divisi"`
	KaryawanUnit   string `json:"karyawan_unit"`
	KaryawanRole   string `json:"karyawan_role"`
	SisaCuti       int    `json:"sisa_cuti"`
}