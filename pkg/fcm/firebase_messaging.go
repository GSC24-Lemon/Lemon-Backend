package fcm

import (
	"context"

	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type FirebaseMessaging struct {
	Client *messaging.Client
}

func NewFirebaseMessaging(projectId string, serviceAccKey string) (*FirebaseMessaging, error) {
	ctx := context.Background()
	// conf := &firebase.Config{ProjectID: projectId}
	opt := option.WithCredentialsFile(serviceAccKey)
	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(ctx)

	if err != nil {
		return nil, err
	}

	fs := &FirebaseMessaging{
		Client: client,
	}

	return fs, nil
}

// defer client.Close()
