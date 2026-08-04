package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type AuthUsecase struct {
	KaryawanRepo     *repository.KaryawanRepo
	RiwayatRepo      *repository.RiwayatRepo
	FCMTokenRepo     *repository.FCMTokenRepo
	SupabaseURL      string
	AnonKey          string
	ResetPasswordURL string
}

func NewAuthUsecase(karyawanRepo *repository.KaryawanRepo, riwayatRepo *repository.RiwayatRepo) *AuthUsecase {
	resetURL := os.Getenv("RESET_PASSWORD_URL")
	if resetURL == "" {
		resetURL = "sakti://reset-password"
	}

	return &AuthUsecase{
		KaryawanRepo:     karyawanRepo,
		RiwayatRepo:      riwayatRepo,
		SupabaseURL:      os.Getenv("SUPABASE_URL"),
		AnonKey:          os.Getenv("SUPABASE_ANON_KEY"),
		ResetPasswordURL: resetURL,
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

	if u.RiwayatRepo != nil {
		detail := "Login berhasil pada " + time.Now().Format("2006-01-02 15:04:05")
		u.RiwayatRepo.CreateRiwayat(ctx, karyawan.ID, "login", detail)
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

func (u *AuthUsecase) ChangePassword(ctx context.Context, userID, token, currentPassword, newPassword string) error {
	log.Printf("ChangePassword: userID=%s", userID)

	karyawan, err := u.KaryawanRepo.GetByID(ctx, userID)
	if err != nil || karyawan == nil {
		return errors.New("karyawan tidak ditemukan")
	}

	verifyReq := map[string]string{
		"email":    karyawan.Email,
		"password": currentPassword,
	}
	jsonBody, _ := json.Marshal(verifyReq)

	verifyHTTPReq, err := http.NewRequestWithContext(ctx, "POST", u.SupabaseURL+"/auth/v1/token?grant_type=password", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	verifyHTTPReq.Header.Set("apikey", u.AnonKey)
	verifyHTTPReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	verifyResp, err := client.Do(verifyHTTPReq)
	if err != nil {
		return err
	}
	defer verifyResp.Body.Close()

	if verifyResp.StatusCode != 200 {
		log.Printf("ChangePassword: current password verification failed, status=%d", verifyResp.StatusCode)
		return errors.New("password saat ini salah")
	}

	changeReq := map[string]string{
		"password": newPassword,
	}
	changeBody, _ := json.Marshal(changeReq)

	changeHTTPReq, err := http.NewRequestWithContext(ctx, "PUT", u.SupabaseURL+"/auth/v1/user", bytes.NewBuffer(changeBody))
	if err != nil {
		return err
	}
	changeHTTPReq.Header.Set("apikey", u.AnonKey)
	changeHTTPReq.Header.Set("Authorization", "Bearer "+token)
	changeHTTPReq.Header.Set("Content-Type", "application/json")

	changeResp, err := client.Do(changeHTTPReq)
	if err != nil {
		return err
	}
	defer changeResp.Body.Close()

	if changeResp.StatusCode != 200 {
		log.Printf("ChangePassword failed: status=%d", changeResp.StatusCode)
		return errors.New("gagal mengubah password")
	}

	if u.RiwayatRepo != nil {
		detail := "Password diubah pada " + time.Now().Format("2006-01-02 15:04:05")
		u.RiwayatRepo.CreateRiwayat(ctx, userID, "update_password", detail)
	}

	log.Printf("ChangePassword berhasil: userID=%s", userID)

	return nil
}

func (u *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	log.Printf("ForgotPassword: email=%s, redirectTo=%s", email, u.ResetPasswordURL)

	karyawan, err := u.KaryawanRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("GetByEmail error: %v", err)
		return errors.New("email tidak ditemukan")
	}
	if karyawan == nil {
		log.Printf("Email %s tidak ditemukan di tabel karyawan", email)
		return errors.New("email tidak ditemukan")
	}

	recoverURL := fmt.Sprintf(
		"%s/auth/v1/recover?redirect_to=%s",
		u.SupabaseURL,
		url.QueryEscape(u.ResetPasswordURL),
	)

	log.Printf("Recover URL: %s", recoverURL)

	supabaseReq := map[string]string{
		"email": email,
	}
	jsonBody, _ := json.Marshal(supabaseReq)

	log.Printf("Supabase payload: %s", string(jsonBody))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", recoverURL, bytes.NewBuffer(jsonBody))
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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read response body error: %v", err)
	}

	log.Printf("Recover response status: %d", resp.StatusCode)
	log.Printf("Recover response body: %s", string(responseBody))

	if resp.StatusCode != 200 {
		log.Printf("ForgotPassword failed: status=%d, body=%s", resp.StatusCode, string(responseBody))
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

	responseBody, _ := io.ReadAll(resp.Body)

	log.Printf("ResetPassword response status: %d", resp.StatusCode)
	log.Printf("ResetPassword response body: %s", string(responseBody))

	if resp.StatusCode != 200 {
		log.Printf("ResetPassword failed: status=%d, body=%s", resp.StatusCode, string(responseBody))
		return errors.New("token tidak valid atau sudah kadaluarsa")
	}

	log.Printf("ResetPassword berhasil")

	return nil
}

func (u *AuthUsecase) Logout(ctx context.Context, token string) error {
	log.Printf("Logout: token=%s...", token[:20])

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u.SupabaseURL+"/auth/v1/logout", nil)
	if err != nil {
		return err
	}
	httpReq.Header.Set("apikey", u.AnonKey)
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("Logout response status: %d", resp.StatusCode)
	log.Printf("Logout response body: %s", string(responseBody))

	if resp.StatusCode != 200 {
		log.Printf("Logout failed: status=%d", resp.StatusCode)
		return errors.New("gagal logout")
	}

	log.Printf("Logout berhasil")
	return nil
}

func (u *AuthUsecase) DeactivateFCMToken(ctx context.Context, karyawanID, fcmToken string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByKaryawanID(ctx, karyawanID)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if token == fcmToken {
			return u.FCMTokenRepo.DeactivateToken(ctx, fcmToken)
		}
	}

	return nil
}