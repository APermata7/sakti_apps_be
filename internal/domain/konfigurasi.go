package domain

import "time"

type KonfigurasiKerja struct {
	ID               string     `json:"id"`
	NamaKantor       string     `json:"nama_kantor"`
	KantorID         *string    `json:"kantor_id"`
	LatKantor        float64    `json:"lat_kantor"`
	LongKantor       float64    `json:"long_kantor"`
	LogoKantor       *string    `json:"logo_kantor"`
	JamMasuk         string     `json:"jam_masuk"`
	JamMinimalMasuk  string     `json:"jam_minimal_masuk"`
	JamPulang        string     `json:"jam_pulang"`
	JamMinimalPulang string     `json:"jam_minimal_pulang"`
	RadiusKantor     int        `json:"radius_kantor"`
	DiperbaruiOleh   *string    `json:"diperbarui_oleh"`
	DiperbaruiPada   time.Time  `json:"diperbarui_pada"`
}

type UpdateKonfigurasiRequest struct {
	NamaKantor       string  `json:"nama_kantor"`
	LatKantor        float64 `json:"lat_kantor"`
	LongKantor       float64 `json:"long_kantor"`
	LogoKantor       *string `json:"logo_kantor"`
	JamMasuk         string  `json:"jam_masuk"`
	JamMinimalMasuk  string  `json:"jam_minimal_masuk"`
	JamPulang        string  `json:"jam_pulang"`
	JamMinimalPulang string  `json:"jam_minimal_pulang"`
	RadiusKantor     int     `json:"radius_kantor"`
}