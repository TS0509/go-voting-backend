package utils

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// InitFirebaseDB 初始化 Firebase Realtime Database 客户端
func InitFirebaseDB() (*db.Client, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("firebase-service-account.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Database(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
