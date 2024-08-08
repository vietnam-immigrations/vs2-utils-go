package push

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func sendToIOSApp(ctx context.Context, title, body string) error {
	log := logger.FromContext(ctx)
	credentials, err := secretsmanager.GetSecret(ctx, "firebaseServiceAccount")
	if err != nil {
		log.Errorf("error getting secret: %v\n", err)
		return err
	}

	// Initialize the Firebase app
	opt := option.WithCredentialsJSON([]byte(credentials))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Errorf("error initializing app: %v\n", err)
		return err
	}

	// Initialize Firebase Cloud Messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Errorf("error getting Messaging client: %v\n", err)
		return err
	}

	// Topic to which you want to send a notification
	topic := "newOrders"

	// Create the message to be sent
	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
				},
			},
		},
	}

	// Send the message
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Errorf("error sending message: %v\n", err)
		return err
	}

	// Response is a message ID string
	log.Infof("Successfully sent message: %s\n", response)
	return nil
}
