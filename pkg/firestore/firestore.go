

package firestore

import (
	"context"
	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

type Firestore struct {
	Client *firestore.Client
}

func NewFirestore(projectId string, serviceAccKey string) (*Firestore, error) {
	ctx := context.Background()
	// conf := &firebase.Config{ProjectID: projectId}
	opt := option.WithCredentialsFile(serviceAccKey)
	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		return nil, err
	}

	fs := &Firestore{
		Client: client,
	}

	return fs, nil
}

// defer client.Close()
