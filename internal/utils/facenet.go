package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

func VerifyFace(ctx context.Context, referenceURL, selfieURL string) (bool, float64, error) {
	faceServiceURL := os.Getenv("FACE_SERVICE_URL")
	if faceServiceURL == "" {
		log.Println("[VerifyFace] FACE_SERVICE_URL not set, bypassing face verification")
		return true, 0.95, nil
	}

	log.Printf("[VerifyFace] Calling face service: reference=%s, selfie=%s", referenceURL, selfieURL)
	return callFaceService(ctx, faceServiceURL, referenceURL, selfieURL)
}

func callFaceService(ctx context.Context, serviceURL, referenceURL, selfieURL string) (bool, float64, error) {
	serviceURL = strings.TrimSuffix(serviceURL, "/")
	url := serviceURL + "/api/v1/verify"

	if referenceURL == "" {
		return false, 0, errors.New("reference URL is empty")
	}
	if selfieURL == "" {
		return false, 0, errors.New("selfie URL is empty")
	}

	reqBody := FaceVerificationRequest{
		ReferenceURL: referenceURL,
		SelfieURL:    selfieURL,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[callFaceService] marshal error: %v", err)
		return false, 0, err
	}

	log.Printf("[callFaceService] POST %s", url)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("[callFaceService] create request error: %v", err)
		return false, 0, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("[callFaceService] request error: %v", err)
		return false, 0, fmt.Errorf("gagal terhubung ke server verifikasi wajah: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[callFaceService] read body error: %v", err)
		return false, 0, err
	}

	log.Printf("[callFaceService] response status=%d, body=%s", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		return false, 0, fmt.Errorf("server verifikasi wajah error (status %d): %s", resp.StatusCode, string(body))
	}

	var result FaceVerificationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[callFaceService] parse error: %v", err)
		return false, 0, err
	}

	if !result.Success {
		log.Printf("[callFaceService] verification failed: %s", result.Message)
		return false, result.Similarity, errors.New(result.Message)
	}

	log.Printf("[callFaceService] verification success: match=%v, similarity=%.4f", result.Match, result.Similarity)
	return result.Match, result.Similarity, nil
}