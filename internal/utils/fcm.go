package utils

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var FCMClient *messaging.Client

func InitFCM() error {
	credPath := os.Getenv("FCM_SERVICE_ACCOUNT")
	if credPath == "" {
		credPath = "firebase-service-account.json"
	}

	if _, err := os.Stat(credPath); os.IsNotExist(err) {
		log.Printf("File service account tidak ditemukan: %s", credPath)
		return nil
	}

	opt := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	FCMClient = client
	log.Println("FCM initialized successfully")
	return nil
}

func SendNotification(token, title, body string) error {
	if FCMClient == nil {
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

	_, err := FCMClient.Send(context.Background(), message)
	if err != nil {
		log.Printf("FCM send notification error: %v", err)
		return err
	}

	log.Printf("FCM notification sent to token: %s", token)
	return nil
}

func SendNotificationWithData(token, title, body string, data map[string]string) error {
	if FCMClient == nil {
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

	_, err := FCMClient.Send(context.Background(), message)
	if err != nil {
		log.Printf("FCM send notification with data error: %v", err)
		return err
	}

	log.Printf("FCM notification with data sent to token: %s", token)
	return nil
}

func SendMulticast(tokens []string, title, body string) error {
	if FCMClient == nil {
		log.Println("FCMClient is nil, cannot send multicast notification")
		return nil
	}

	if len(tokens) == 0 {
		log.Println("No tokens to send multicast notification")
		return nil
	}

	log.Printf("Sending FCM multicast notification to %d devices", len(tokens))

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: tokens,
	}

	response, err := FCMClient.SendMulticast(context.Background(), message)
	if err != nil {
		log.Printf("FCM send multicast error: %v", err)
		return err
	}

	log.Printf("FCM successCount: %d, failureCount: %d", response.SuccessCount, response.FailureCount)

	if response.FailureCount > 0 {
		for i, result := range response.Responses {
			if result.Error != nil {
				log.Printf("FCM failure for token index %d: %v", i, result.Error)
			}
		}
	}

	return nil
}

func SendMulticastWithData(tokens []string, title, body string, data map[string]string) error {
	if FCMClient == nil {
		log.Println("FCMClient is nil, cannot send multicast with data")
		return nil
	}

	if len(tokens) == 0 {
		log.Println("No tokens to send multicast with data")
		return nil
	}

	log.Printf("Sending FCM multicast with data to %d devices", len(tokens))

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:   data,
		Tokens: tokens,
	}

	response, err := FCMClient.SendMulticast(context.Background(), message)
	if err != nil {
		log.Printf("FCM send multicast with data error: %v", err)
		return err
	}

	log.Printf("FCM successCount: %d, failureCount: %d", response.SuccessCount, response.FailureCount)

	if response.FailureCount > 0 {
		for i, result := range response.Responses {
			if result.Error != nil {
				log.Printf("FCM failure for token index %d: %v", i, result.Error)
			}
		}
	}

	return nil
}

func IsFCMReady() bool {
	return FCMClient != nil
}