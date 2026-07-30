package utils

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMService struct {
	client *messaging.Client
	mu     sync.RWMutex
	ready  bool
}

var (
	fcmService *FCMService
	serviceMu  sync.Mutex
)

func GetFCMService() *FCMService {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	if fcmService == nil {
		fcmService = &FCMService{}
	}
	return fcmService
}

func (s *FCMService) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ready && s.client != nil {
		log.Println("FCM already initialized")
		return nil
	}

	log.Println("Initializing FCM...")

	credBase64 := os.Getenv("FCM_SERVICE_ACCOUNT_BASE64")
	credPath := os.Getenv("FCM_SERVICE_ACCOUNT")

	var opt option.ClientOption

	if credBase64 != "" {
		log.Println("Using FCM_SERVICE_ACCOUNT_BASE64")
		decoded, err := base64.StdEncoding.DecodeString(credBase64)
		if err != nil {
			log.Printf("Failed to decode base64: %v", err)
			return err
		}
		opt = option.WithCredentialsJSON(decoded)
	} else if credPath != "" {
		log.Printf("Using FCM_SERVICE_ACCOUNT file: %s", credPath)
		if _, err := os.Stat(credPath); os.IsNotExist(err) {
			log.Printf("File not found: %s", credPath)
			return err
		}
		opt = option.WithCredentialsFile(credPath)
	} else {
		credPath = "firebase-service-account.json"
		if _, err := os.Stat(credPath); os.IsNotExist(err) {
			log.Printf("No service account found")
			return nil
		}
		opt = option.WithCredentialsFile(credPath)
	}

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Firebase NewApp error: %v", err)
		return err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Printf("App Messaging error: %v", err)
		return err
	}

	s.client = client
	s.ready = true
	log.Println("FCM initialized successfully")
	return nil
}

func (s *FCMService) ensureClient() bool {
	s.mu.RLock()
	if s.client != nil && s.ready {
		s.mu.RUnlock()
		return true
	}
	s.mu.RUnlock()

	log.Println("FCMClient is nil or not ready, attempting to reinitialize...")
	if err := s.Init(); err != nil {
		log.Printf("Failed to reinitialize FCM: %v", err)
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.client == nil {
		log.Println("FCMClient still nil after reinitialization")
		return false
	}

	log.Println("FCMClient reinitialized successfully")
	return true
}

func (s *FCMService) GetClient() *messaging.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.client
}

func (s *FCMService) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ready && s.client != nil
}

func InitFCM() error {
	return GetFCMService().Init()
}

func SendNotification(token, title, body string) error {
	service := GetFCMService()
	if !service.ensureClient() {
		return nil
	}

	client := service.GetClient()
	if client == nil {
		log.Println("FCMClient is nil, cannot send notification")
		return nil
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}

	_, err := client.Send(context.Background(), message)
	if err != nil {
		log.Printf("FCM send notification error: %v", err)
		return err
	}

	log.Printf("FCM notification sent to token: %s", token)
	return nil
}

func SendNotificationWithData(token, title, body string, data map[string]string) error {
	service := GetFCMService()
	if !service.ensureClient() {
		return nil
	}

	client := service.GetClient()
	if client == nil {
		log.Println("FCMClient is nil, cannot send notification with data")
		return nil
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:  data,
		Token: token,
	}

	_, err := client.Send(context.Background(), message)
	if err != nil {
		log.Printf("FCM send notification with data error: %v", err)
		return err
	}

	log.Printf("FCM notification with data sent to token: %s", token)
	return nil
}

func SendMulticast(tokens []string, title, body string) error {
	service := GetFCMService()
	if !service.ensureClient() {
		return nil
	}

	if len(tokens) == 0 {
		log.Println("No tokens to send multicast notification")
		return nil
	}

	client := service.GetClient()
	if client == nil {
		log.Println("FCMClient is nil, cannot send multicast notification")
		return nil
	}

	log.Printf("Sending FCM multicast notification to %d devices", len(tokens))

	successCount := 0
	failureCount := 0

	for i, token := range tokens {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Token: token,
		}

		_, err := client.Send(context.Background(), message)
		if err != nil {
			log.Printf("FCM send error for token %d: %v", i, err)
			failureCount++
		} else {
			successCount++
		}
	}

	log.Printf("FCM successCount: %d, failureCount: %d", successCount, failureCount)

	if failureCount > 0 {
		log.Printf("FCM some notifications failed to send")
	}

	return nil
}

func SendMulticastWithData(tokens []string, title, body string, data map[string]string) error {
	service := GetFCMService()
	if !service.ensureClient() {
		return nil
	}

	if len(tokens) == 0 {
		log.Println("No tokens to send multicast with data")
		return nil
	}

	client := service.GetClient()
	if client == nil {
		log.Println("FCMClient is nil, cannot send multicast with data")
		return nil
	}

	log.Printf("Sending FCM multicast with data to %d devices", len(tokens))

	successCount := 0
	failureCount := 0

	for i, token := range tokens {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Data:  data,
			Token: token,
		}

		_, err := client.Send(context.Background(), message)
		if err != nil {
			log.Printf("FCM send with data error for token %d: %v", i, err)
			failureCount++
		} else {
			successCount++
		}
	}

	log.Printf("FCM successCount: %d, failureCount: %d", successCount, failureCount)

	if failureCount > 0 {
		log.Printf("FCM some notifications with data failed to send")
	}

	return nil
}

func IsFCMReady() bool {
	return GetFCMService().IsReady()
}