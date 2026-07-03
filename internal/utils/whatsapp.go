package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WhatsAppMessage struct {
	MessagingProduct string              `json:"messaging_product"`
	To               string              `json:"to"`
	Type             string              `json:"type"`
	Template         WhatsAppTemplate    `json:"template"`
}

type WhatsAppTemplate struct {
	Name       string              `json:"name"`
	Language   WhatsAppLanguage    `json:"language"`
	Components []WhatsAppComponent `json:"components"`
}

type WhatsAppLanguage struct {
	Code string `json:"code"`
}

type WhatsAppComponent struct {
	Type       string              `json:"type"`
	Parameters []WhatsAppParameter `json:"parameters"`
}

type WhatsAppParameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func SendWhatsAppNotification(to, namaAtasan, namaPegawai, jenis, tanggal string) error {
	phoneNumberID := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	accessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN")

	if phoneNumberID == "" || accessToken == "" {
		return nil
	}

	message := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: WhatsAppTemplate{
			Name:     "leave_notification",
			Language: WhatsAppLanguage{Code: "id"},
			Components: []WhatsAppComponent{
				{
					Type: "body",
					Parameters: []WhatsAppParameter{
						{Type: "text", Text: namaAtasan},
						{Type: "text", Text: namaPegawai},
						{Type: "text", Text: jenis},
						{Type: "text", Text: tanggal},
					},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(message)
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", phoneNumberID)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("whatsapp error: %d", resp.StatusCode)
	}

	return nil
}