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
	return err
}

func SendNotificationWithData(token, title, body string, data map[string]string) error {
	if FCMClient == nil {
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
	return err
}

func SendMulticast(tokens []string, title, body string) error {
	if FCMClient == nil || len(tokens) == 0 {
		return nil
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: tokens,
	}

	_, err := FCMClient.SendMulticast(context.Background(), message)
	return err
}

func SendMulticastWithData(tokens []string, title, body string, data map[string]string) error {
	if FCMClient == nil || len(tokens) == 0 {
		return nil
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:   data,
		Tokens: tokens,
	}

	_, err := FCMClient.SendMulticast(context.Background(), message)
	return err
}

func IsFCMReady() bool {
	return FCMClient != nil
}