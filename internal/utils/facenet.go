package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"sakti_apps_be/internal/repository"
)

type FaceVerificationResponse struct {
	Match      bool    `json:"match"`
	Similarity float64 `json:"similarity"`
}

type FaceVerificationRequest struct {
	ReferenceURL string `json:"reference_url"`
	SelfieURL    string `json:"selfie_url"`
}

func VerifyFace(ctx context.Context, selfieURL, employeeID string) (bool, float64, error) {
	faceServiceURL := os.Getenv("FACE_SERVICE_URL")
	if faceServiceURL == "" {
		return true, 0.95, nil
	}

	return true, 0.95, nil
}

func VerifyFaceWithRepo(ctx context.Context, selfieURL, employeeID string, karyawanRepo *repository.KaryawanRepo) (bool, float64, error) {
	faceServiceURL := os.Getenv("FACE_SERVICE_URL")
	if faceServiceURL == "" {
		return true, 0.95, nil
	}

	karyawan, err := karyawanRepo.GetByID(ctx, employeeID)
	if err != nil || karyawan == nil {
		return false, 0, errors.New("karyawan tidak ditemukan")
	}

	if karyawan.FotoURL == nil {
		return false, 0, errors.New("foto wajah tidak ditemukan")
	}

	reqBody := FaceVerificationRequest{
		ReferenceURL: *karyawan.FotoURL,
		SelfieURL:    selfieURL,
	}
	jsonBody, _ := json.Marshal(reqBody)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", faceServiceURL+"/verify", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, 0, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, err
	}

	var result FaceVerificationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return false, 0, err
	}

	return result.Match, result.Similarity, nil
}