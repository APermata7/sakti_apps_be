package usecase

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "errors"

    "sakti_apps_be/internal/domain"
    "sakti_apps_be/internal/repository"
)

type TTDUsecase struct {
    TTDRepo *repository.TTDRepo
}

func NewTTDUsecase(ttdRepo *repository.TTDRepo) *TTDUsecase {
    return &TTDUsecase{TTDRepo: ttdRepo}
}

func (u *TTDUsecase) Create(ctx context.Context, karyawanID string, req domain.CreateTTDRequest) (*domain.TandaTangan, error) {
    if req.URLTandaTangan == "" {
        return nil, errors.New("URL tanda tangan wajib diisi")
    }

    hash := sha256.Sum256([]byte(req.URLTandaTangan))
    hashString := hex.EncodeToString(hash[:])

    ttd := &domain.TandaTangan{
        KaryawanID:      karyawanID,
        URLTandaTangan:  req.URLTandaTangan,
        HashTandaTangan: &hashString,
    }

    err := u.TTDRepo.Create(ctx, ttd)
    if err != nil {
        return nil, err
    }

    return ttd, nil
}

func (u *TTDUsecase) GetByKaryawanID(ctx context.Context, karyawanID string) (*domain.TandaTangan, error) {
    return u.TTDRepo.GetByKaryawanID(ctx, karyawanID)
}

func (u *TTDUsecase) Update(ctx context.Context, karyawanID string, req domain.CreateTTDRequest) (*domain.TandaTangan, error) {
    if req.URLTandaTangan == "" {
        return nil, errors.New("URL tanda tangan wajib diisi")
    }

    hash := sha256.Sum256([]byte(req.URLTandaTangan))
    hashString := hex.EncodeToString(hash[:])

    ttd := &domain.TandaTangan{
        KaryawanID:      karyawanID,
        URLTandaTangan:  req.URLTandaTangan,
        HashTandaTangan: &hashString,
    }

    err := u.TTDRepo.Update(ctx, ttd)
    if err != nil {
        return nil, err
    }

    return ttd, nil
}

func (u *TTDUsecase) Delete(ctx context.Context, karyawanID string) error {
    return u.TTDRepo.Delete(ctx, karyawanID)
}