package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type AuthUsecase struct {
	KaryawanRepo *repository.KaryawanRepo
	SupabaseURL  string
	AnonKey      string
}

func NewAuthUsecase(karyawanRepo *repository.KaryawanRepo) *AuthUsecase {
	return &AuthUsecase{
		KaryawanRepo: karyawanRepo,
		SupabaseURL:  os.Getenv("SUPABASE_URL"),
		AnonKey:      os.Getenv("SUPABASE_ANON_KEY"),
	}
}

func (u *AuthUsecase) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	supabaseReq := map[string]string{
		"email":    req.Email,
		"password": req.Password,
	}
	jsonBody, _ := json.Marshal(supabaseReq)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u.SupabaseURL+"/auth/v1/token?grant_type=password", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("apikey", u.AnonKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New("email atau password salah")
	}

	var authResp struct {
		AccessToken string `json:"access_token"`
		User        struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}
	if err := json.Unmarshal(body, &authResp); err != nil {
		return nil, err
	}

	karyawan, err := u.KaryawanRepo.GetByID(ctx, authResp.User.ID)
	if err != nil {
		return nil, err
	}

	if karyawan == nil {
		return nil, errors.New("data karyawan tidak ditemukan")
	}

	if karyawan.StatusKaryawan != "aktif" {
		return nil, errors.New("akun tidak aktif")
	}

	return &domain.LoginResponse{
		AccessToken: authResp.AccessToken,
		User:        *karyawan,
	}, nil
}

func (u *AuthUsecase) GetProfile(ctx context.Context, userID string) (*domain.Karyawan, error) {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if karyawan == nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}

	return karyawan, nil
}

func (u *AuthUsecase) CreateKaryawan(ctx context.Context, req domain.CreateKaryawanRequest) error {
	return nil
}

func (u *AuthUsecase) UpdateKaryawan(ctx context.Context, id string, req domain.UpdateKaryawanRequest) error {
	return nil
}

func (u *AuthUsecase) DeleteKaryawan(ctx context.Context, id string) error {
	return nil
}