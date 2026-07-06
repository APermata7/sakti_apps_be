package usecase

import (
	"context"
	"errors"
	"time"

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

	libur := &domain.Libur{
		Tanggal: t.Format("2006-01-02"),
		Nama:    req.Nama,
		Jenis:   req.Jenis,
		Aktif:   true,
		Sumber:  "manual",
	}

	if err := u.LiburRepo.Create(ctx, libur); err != nil {
		return nil, err
	}
	return libur, nil
}

func (u *LiburUsecase) GetAll(ctx context.Context, tahun string) ([]domain.Libur, error) {
	if tahun == "" {
		tahun = time.Now().Format("2006")
	}
	return u.LiburRepo.GetAll(ctx, tahun)
}

func (u *LiburUsecase) GetByID(ctx context.Context, id string) (*domain.Libur, error) {
	return u.LiburRepo.GetByID(ctx, id)
}

func (u *LiburUsecase) Update(ctx context.Context, id string, req domain.UpdateLiburRequest) (*domain.Libur, error) {
	libur, err := u.LiburRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if libur == nil {
		return nil, errors.New("libur tidak ditemukan")
	}

	if req.Nama != "" {
		libur.Nama = req.Nama
	}
	if req.Jenis != "" {
		libur.Jenis = req.Jenis
	}
	libur.Aktif = req.Aktif

	if err := u.LiburRepo.Update(ctx, libur); err != nil {
		return nil, err
	}
	return libur, nil
}

func (u *LiburUsecase) Delete(ctx context.Context, id string) error {
	return u.LiburRepo.Delete(ctx, id)
}

func (u *LiburUsecase) IsHoliday(ctx context.Context, tanggal string) (bool, error) {
	return u.LiburRepo.IsHoliday(ctx, tanggal)
}