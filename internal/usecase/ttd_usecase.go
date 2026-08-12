package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"time"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type TTDUsecase struct {
	TTDRepo *repository.TTDRepo
}

func NewTTDUsecase(ttdRepo *repository.TTDRepo) *TTDUsecase {
	return &TTDUsecase{TTDRepo: ttdRepo}
}

func (u *TTDUsecase) generateHashFromURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("url tidak boleh kosong")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("gagal mengunduh file tanda tangan")
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, resp.Body); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (u *TTDUsecase) Create(ctx context.Context, karyawanID string, req domain.CreateTTDRequest) (*domain.TandaTangan, error) {
	if req.URLTandaTangan == "" {
		return nil, errors.New("url tanda tangan wajib diisi")
	}

	existing, err := u.TTDRepo.GetByKaryawanID(ctx, karyawanID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("karyawan sudah memiliki tanda tangan, gunakan endpoint update")
	}

	hashString, err := u.generateHashFromURL(req.URLTandaTangan)
	if err != nil {
		return nil, errors.New("gagal memproses hash tanda tangan")
	}

	ttd := &domain.TandaTangan{
		KaryawanID:      karyawanID,
		URLTandaTangan:  req.URLTandaTangan,
		HashTandaTangan: &hashString,
	}

	if err := u.TTDRepo.Create(ctx, ttd); err != nil {
		return nil, err
	}

	return ttd, nil
}

func (u *TTDUsecase) GetByKaryawanID(ctx context.Context, karyawanID string) (*domain.TandaTangan, error) {
	return u.TTDRepo.GetByKaryawanID(ctx, karyawanID)
}

func (u *TTDUsecase) Update(ctx context.Context, karyawanID string, req domain.CreateTTDRequest) (*domain.TandaTangan, error) {
	if req.URLTandaTangan == "" {
		return nil, errors.New("url tanda tangan wajib diisi")
	}

	existing, err := u.TTDRepo.GetByKaryawanID(ctx, karyawanID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("karyawan belum memiliki tanda tangan")
	}

	hashString, err := u.generateHashFromURL(req.URLTandaTangan)
	if err != nil {
		return nil, errors.New("gagal memproses hash tanda tangan")
	}

	existing.URLTandaTangan = req.URLTandaTangan
	existing.HashTandaTangan = &hashString

	if err := u.TTDRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (u *TTDUsecase) Delete(ctx context.Context, karyawanID string) error {
	existing, err := u.TTDRepo.GetByKaryawanID(ctx, karyawanID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("karyawan belum memiliki tanda tangan")
	}
	return u.TTDRepo.Delete(ctx, karyawanID)
}

func (u *TTDUsecase) Verify(ctx context.Context, karyawanID string) (*domain.VerifyTTDResponse, error) {
	ttd, err := u.TTDRepo.GetByKaryawanID(ctx, karyawanID)
	if err != nil {
		return nil, err
	}
	if ttd == nil {
		return &domain.VerifyTTDResponse{
			IsValid: false,
			Message: "tanda tangan tidak ditemukan",
		}, nil
	}

	if ttd.HashTandaTangan == nil {
		return &domain.VerifyTTDResponse{
			IsValid: false,
			Message: "hash tanda tangan tidak tersedia",
		}, nil
	}

	currentHash, err := u.generateHashFromURL(ttd.URLTandaTangan)
	if err != nil {
		return &domain.VerifyTTDResponse{
			IsValid: false,
			Message: "gagal verifikasi tanda tangan: " + err.Error(),
		}, nil
	}

	isValid := currentHash == *ttd.HashTandaTangan
	message := "tanda tangan valid"
	if !isValid {
		message = "tanda tangan tidak valid, file telah diubah"
	}

	return &domain.VerifyTTDResponse{
		IsValid:   isValid,
		Message:   message,
		UpdatedAt: ttd.DiperbaruiPada.Format("2006-01-02 15:04:05"),
	}, nil
}