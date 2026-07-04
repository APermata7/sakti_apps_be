package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type KonfigurasiRepo struct {
	DB *pgxpool.Pool
}

func NewKonfigurasiRepo(db *pgxpool.Pool) *KonfigurasiRepo {
	return &KonfigurasiRepo{DB: db}
}

func (r *KonfigurasiRepo) GetActive(ctx context.Context) (*domain.KonfigurasiKerja, error) {
	query := `
		SELECT id, nama_kantor, lat_kantor, long_kantor, logo_kantor,
		       jam_masuk, jam_minimal_masuk, jam_pulang, jam_minimal_pulang,
		       radius_kantor, diperbarui_oleh, diperbarui_pada
		FROM konfigurasi_kerja
		LIMIT 1
	`

	var k domain.KonfigurasiKerja
	err := r.DB.QueryRow(ctx, query).Scan(
		&k.ID, &k.NamaKantor, &k.LatKantor, &k.LongKantor, &k.LogoKantor,
		&k.JamMasuk, &k.JamMinimalMasuk, &k.JamPulang, &k.JamMinimalPulang,
		&k.RadiusKantor, &k.DiperbaruiOleh, &k.DiperbaruiPada,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &k, nil
}

func (r *KonfigurasiRepo) Update(ctx context.Context, k *domain.KonfigurasiKerja) error {
	query := `
		UPDATE konfigurasi_kerja 
		SET nama_kantor = $2, lat_kantor = $3, long_kantor = $4, logo_kantor = $5,
		    jam_masuk = $6, jam_minimal_masuk = $7, jam_pulang = $8, jam_minimal_pulang = $9,
		    radius_kantor = $10, diperbarui_oleh = $11, diperbarui_pada = NOW()
		WHERE id = $1
	`
	_, err := r.DB.Exec(ctx, query,
		k.ID, k.NamaKantor, k.LatKantor, k.LongKantor, k.LogoKantor,
		k.JamMasuk, k.JamMinimalMasuk, k.JamPulang, k.JamMinimalPulang,
		k.RadiusKantor, k.DiperbaruiOleh,
	)
	return err
}