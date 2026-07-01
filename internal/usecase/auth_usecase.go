package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
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

func (u *AuthUsecase) ChangePassword(ctx context.Context, userID, newPassword string) error {
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
		return errors.New("gagal mengubah password")
	}

	return nil
}

func (u *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
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
		return errors.New("gagal mengirim link reset")
	}

	return nil
}

func (u *AuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
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
		return errors.New("token tidak valid atau sudah kadaluarsa")
	}

	return nil
}

func (u *AuthUsecase) CreateKaryawan(ctx context.Context, req domain.CreateKaryawanRequest) error {
	supabaseReq := map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
		"user_metadata": map[string]string{
			"nama_lengkap":  req.NamaLengkap,
			"peran":         req.Peran,
			"level_jabatan": req.LevelJabatan,
		},
	}
	jsonBody, _ := json.Marshal(supabaseReq)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u.SupabaseURL+"/auth/v1/signup", bytes.NewBuffer(jsonBody))
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
		return errors.New("gagal membuat akun")
	}

	var authResp struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	karyawan := &domain.Karyawan{
		ID:           authResp.ID,
		NamaLengkap:  req.NamaLengkap,
		Email:        req.Email,
		Peran:        req.Peran,
		LevelJabatan: req.LevelJabatan,
		StatusKaryawan: "aktif",
	}

	if req.NomorTelepon != "" {
		karyawan.NomorTelepon = &req.NomorTelepon
	}
	if req.FotoURL != "" {
		karyawan.FotoURL = &req.FotoURL
	}
	if req.Divisi != "" {
		karyawan.Divisi = &req.Divisi
	}
	if req.Unit != "" {
		karyawan.Unit = &req.Unit
	}
	if req.AtasanLangsungID != "" {
		karyawan.AtasanLangsungID = &req.AtasanLangsungID
	}

	return u.KaryawanRepo.Create(ctx, karyawan)
}

func (u *AuthUsecase) UpdateKaryawan(ctx context.Context, id string, req domain.UpdateKaryawanRequest) error {
	existing, err := u.KaryawanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("karyawan tidak ditemukan")
	}

	if req.NamaLengkap != "" {
		existing.NamaLengkap = req.NamaLengkap
	}
	if req.NomorTelepon != "" {
		existing.NomorTelepon = &req.NomorTelepon
	}
	if req.FotoURL != "" {
		existing.FotoURL = &req.FotoURL
	}
	if req.Peran != "" {
		existing.Peran = req.Peran
	}
	if req.LevelJabatan != "" {
		existing.LevelJabatan = req.LevelJabatan
	}
	if req.AtasanLangsungID != "" {
		existing.AtasanLangsungID = &req.AtasanLangsungID
	}
	if req.Divisi != "" {
		existing.Divisi = &req.Divisi
	}
	if req.Unit != "" {
		existing.Unit = &req.Unit
	}
	if req.StatusKaryawan != "" {
		existing.StatusKaryawan = req.StatusKaryawan
	}

	return u.KaryawanRepo.Update(ctx, existing)
}

func (u *AuthUsecase) DeleteKaryawan(ctx context.Context, id string) error {
	return u.KaryawanRepo.Delete(ctx, id)
}