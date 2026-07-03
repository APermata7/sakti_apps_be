package domain

type RiwayatItem struct {
	Tipe       string `json:"tipe"`
	ID         string `json:"id"`
	Tanggal    string `json:"tanggal"`
	JamMasuk   string `json:"jam_masuk,omitempty"`
	JamKeluar  string `json:"jam_keluar,omitempty"`
	Status     string `json:"status"`
	TotalHari  int    `json:"total_hari,omitempty"`
	Lembur     bool   `json:"lembur,omitempty"`
	JamLembur  float64 `json:"jam_lembur,omitempty"`
	URLFoto    string `json:"url_foto,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type RiwayatResponse struct {
	Items []RiwayatItem `json:"items"`
	Meta  MetaPagination `json:"meta"`
}