// package gorm

// import (
// 	"fmt"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// type Gorm struct {
// 	Pool *gorm.DB
// }

// func NewGorm(username string, password string) (*Gorm, error) {
// 	dsn := "host=localhost user=" + username + " password=" + password + " dbname=chat port=5432 sslmode=disable TimeZone=Asia/Shanghai"
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("gorm - NewGorm - gorm.Open: %w", err)
// 	}
// 	gorm := &Gorm{
// 		Pool: db,
// 	}

// 	return gorm, nil
// }

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
	//conf := &firebase.Config{ProjectID: projectId}
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
