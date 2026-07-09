package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
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

	log.Printf("Login attempt: email=%s", req.Email)

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

	if resp.StatusCode != 200 {
		log.Printf("Supabase login failed: status=%d", resp.StatusCode)
		return nil, errors.New("email atau password salah")
	}

	var authResp struct {
		AccessToken string `json:"access_token"`
		User        struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	log.Printf("Auth User ID: %s, Email: %s", authResp.User.ID, authResp.User.Email)

	karyawan, err := u.KaryawanRepo.GetByID(ctx, authResp.User.ID)
	if err != nil {
		log.Printf("GetByID error: %v", err)
		return nil, err
	}

	if karyawan == nil {
		log.Printf("Karyawan dengan ID %s tidak ditemukan di tabel karyawan", authResp.User.ID)
		return nil, errors.New("data karyawan tidak ditemukan")
	}

	log.Printf("Karyawan ditemukan: ID=%s, Status=%s, Role=%s", karyawan.ID, karyawan.StatusKaryawan, karyawan.Role)

	if karyawan.StatusKaryawan != "aktif" {
		log.Printf("Login ditolak: status karyawan = %s", karyawan.StatusKaryawan)
		return nil, errors.New("akun tidak aktif")
	}

	log.Printf("Login berhasil: email=%s, role=%s", karyawan.Email, karyawan.Role)

	return &domain.LoginResponse{
		AccessToken: authResp.AccessToken,
		User:        *karyawan,
	}, nil
}

func (u *AuthUsecase) GetProfile(ctx context.Context, userID string) (*domain.Karyawan, error) {
	log.Printf("GetProfile: userID=%s", userID)

	karyawan, err := u.KaryawanRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("GetByID error: %v", err)
		return nil, err
	}

	if karyawan == nil {
		log.Printf("Karyawan dengan ID %s tidak ditemukan", userID)
		return nil, errors.New("karyawan tidak ditemukan")
	}

	log.Printf("Profile ditemukan: email=%s, role=%s", karyawan.Email, karyawan.Role)

	return karyawan, nil
}

func (u *AuthUsecase) ChangePassword(ctx context.Context, userID, newPassword string) error {
	log.Printf("ChangePassword: userID=%s", userID)

	supabaseReq := map[string]string{
		"password": newPassword,
	}
	jsonBody, _ := json.Marshal(supabaseReq)

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", u.SupabaseURL+"/auth/v1/user", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("apikey", u.AnonKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("ChangePassword failed: status=%d", resp.StatusCode)
		return errors.New("gagal mengubah password")
	}

	log.Printf("ChangePassword berhasil: userID=%s", userID)

	return nil
}

func (u *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	log.Printf("ForgotPassword: email=%s", email)

	karyawan, err := u.KaryawanRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("GetByEmail error: %v", err)
		return errors.New("email tidak ditemukan")
	}
	if karyawan == nil {
		log.Printf("Email %s tidak ditemukan di tabel karyawan", email)
		return errors.New("email tidak ditemukan")
	}

	supabaseReq := map[string]string{
		"email": email,
	}
	jsonBody, _ := json.Marshal(supabaseReq)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u.SupabaseURL+"/auth/v1/recover", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("apikey", u.AnonKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("ForgotPassword failed: status=%d", resp.StatusCode)
		return errors.New("gagal mengirim link reset")
	}

	log.Printf("ForgotPassword berhasil: email=%s", email)

	return nil
}

func (u *AuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	log.Printf("ResetPassword: token=%s...", token[:20])

	supabaseReq := map[string]string{
		"password": newPassword,
	}
	jsonBody, _ := json.Marshal(supabaseReq)

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", u.SupabaseURL+"/auth/v1/user", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("apikey", u.AnonKey)
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("ResetPassword failed: status=%d", resp.StatusCode)
		return errors.New("token tidak valid atau sudah kadaluarsa")
	}

	log.Printf("ResetPassword berhasil")

	return nil
}