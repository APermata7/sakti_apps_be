package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type KonfigurasiRepo struct {
	DB *pgxpool.Pool
}

func NewKonfigurasiRepo(db *pgxpool.Pool) *KonfigurasiRepo {
	return &KonfigurasiRepo{DB: db}
}

type Konfigurasi struct {
	ID               string  `json:"id"`
	NamaKantor       string  `json:"nama_kantor"`
	LatKantor        float64 `json:"lat_kantor"`
	LongKantor       float64 `json:"long_kantor"`
	JamMasuk         string  `json:"jam_masuk"`
	JamMinimalMasuk  string  `json:"jam_minimal_masuk"`
	JamPulang        string  `json:"jam_pulang"`
	JamMinimalPulang string  `json:"jam_minimal_pulang"`
	RadiusKantor     int     `json:"radius_kantor"`
}

func (r *KonfigurasiRepo) GetActive(ctx context.Context) (*Konfigurasi, error) {
	query := `
		SELECT id, nama_kantor, lat_kantor, long_kantor,
		       jam_masuk, jam_minimal_masuk, jam_pulang, jam_minimal_pulang,
		       radius_kantor
		FROM konfigurasi_kerja
		LIMIT 1
	`

	var k Konfigurasi
	err := r.DB.QueryRow(ctx, query).Scan(
		&k.ID, &k.NamaKantor, &k.LatKantor, &k.LongKantor,
		&k.JamMasuk, &k.JamMinimalMasuk, &k.JamPulang, &k.JamMinimalPulang,
		&k.RadiusKantor,
	)
	if err != nil {
		return nil, err
	}

	return &k, nil
}