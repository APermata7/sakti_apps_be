package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type LiburUsecase struct {
	LiburRepo *repository.LiburRepo
}

func NewLiburUsecase(liburRepo *repository.LiburRepo) *LiburUsecase {
	return &LiburUsecase{LiburRepo: liburRepo}
}

func (u *LiburUsecase) Create(ctx context.Context, req domain.CreateLiburRequest) (*domain.Libur, error) {
	if req.Tanggal == "" || req.Nama == "" || req.Jenis == "" {
		return nil, errors.New("tanggal, nama, dan jenis wajib diisi")
	}

	t, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid (YYYY-MM-DD)")
	}

	if req.Jenis != "nasional" && req.Jenis != "cuti_bersama" {
		return nil, errors.New("jenis harus 'nasional' atau 'cuti_bersama'")
	}

	sumber := req.Sumber
	if sumber == "" {
		sumber = "manual"
	}
	if sumber != "api" && sumber != "manual" {
		return nil, errors.New("sumber harus 'api' atau 'manual'")
	}

	libur := &domain.Libur{
		Tanggal: t,
		Nama:    req.Nama,
		Jenis:   req.Jenis,
		Aktif:   true,
		Sumber:  sumber,
	}

	if err := u.LiburRepo.Create(ctx, libur); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("tanggal sudah terdaftar sebagai hari libur")
		}
		return nil, err
	}
	return libur, nil
}

func (u *LiburUsecase) GetAll(ctx context.Context, tahun int, jenis, sumber string, aktif *bool, limit, page int) ([]domain.Libur, int, error) {
	offset := (page - 1) * limit
	return u.LiburRepo.GetAll(ctx, tahun, jenis, sumber, aktif, limit, offset)
}

func (u *LiburUsecase) GetByID(ctx context.Context, id string) (*domain.Libur, error) {
	libur, err := u.LiburRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if libur == nil {
		return nil, errors.New("hari libur tidak ditemukan")
	}
	return libur, nil
}

func (u *LiburUsecase) GetByTahun(ctx context.Context, tahun int) ([]domain.Libur, error) {
	return u.LiburRepo.GetByTahun(ctx, tahun)
}

func (u *LiburUsecase) GetByBulan(ctx context.Context, bulan string) ([]domain.Libur, error) {
	return u.LiburRepo.GetByBulan(ctx, bulan)
}

func (u *LiburUsecase) CheckTanggal(ctx context.Context, tanggal string) (*domain.Libur, bool, error) {
	libur, err := u.LiburRepo.GetByTanggal(ctx, tanggal)
	if err != nil {
		return nil, false, err
	}
	if libur == nil {
		return nil, false, nil
	}
	return libur, true, nil
}

func (u *LiburUsecase) GetAktif(ctx context.Context) ([]domain.Libur, error) {
	aktif := true
	items, _, err := u.LiburRepo.GetAll(ctx, 0, "", "", &aktif, 100, 0)
	return items, err
}

func (u *LiburUsecase) Update(ctx context.Context, id string, req domain.UpdateLiburRequest) (*domain.Libur, error) {
	libur, err := u.LiburRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if libur == nil {
		return nil, errors.New("hari libur tidak ditemukan")
	}

	if req.Nama != "" {
		libur.Nama = req.Nama
	}
	if req.Jenis != "" {
		if req.Jenis != "nasional" && req.Jenis != "cuti_bersama" {
			return nil, errors.New("jenis harus 'nasional' atau 'cuti_bersama'")
		}
		libur.Jenis = req.Jenis
	}

	err = u.LiburRepo.Update(ctx, libur)
	if err != nil {
		return nil, err
	}
	return libur, nil
}

func (u *LiburUsecase) Delete(ctx context.Context, id string) error {
	libur, err := u.LiburRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if libur == nil {
		return errors.New("hari libur tidak ditemukan")
	}
	return u.LiburRepo.Delete(ctx, id)
}

func (u *LiburUsecase) Toggle(ctx context.Context, id string) (bool, error) {
	return u.LiburRepo.Toggle(ctx, id)
}

func (u *LiburUsecase) IsHoliday(ctx context.Context, tanggal string) (bool, error) {
	libur, err := u.LiburRepo.GetByTanggal(ctx, tanggal)
	if err != nil {
		return false, err
	}
	if libur == nil {
		return false, nil
	}
	return libur.Aktif, nil
}