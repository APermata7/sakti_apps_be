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
	Success    bool    `json:"success"`
	Match      bool    `json:"match"`
	Similarity float64 `json:"similarity"`
	Threshold  float64 `json:"threshold"`
	Message    string  `json:"message"`
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

	return callFaceService(ctx, faceServiceURL, "", selfieURL)
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

	return callFaceService(ctx, faceServiceURL, *karyawan.FotoURL, selfieURL)
}

func callFaceService(ctx context.Context, serviceURL, referenceURL, selfieURL string) (bool, float64, error) {
	reqBody := FaceVerificationRequest{
		ReferenceURL: referenceURL,
		SelfieURL:    selfieURL,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return false, 0, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		serviceURL+"/api/v1/verify", 
		bytes.NewBuffer(jsonBody),
	)
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

	if !result.Success {
		return false, result.Similarity, errors.New(result.Message)
	}

	return result.Match, result.Similarity, nil
}