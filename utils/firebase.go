package utils

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	firestoreClient *firestore.Client
	initOnce        sync.Once
	initErr         error
)

func GetFirestoreClient() (*firestore.Client, error) {
	initOnce.Do(func() {
		ctx := context.Background()
		sa := option.WithCredentialsFile("config/firebase-service-account.json")
		firestoreClient, initErr = firestore.NewClient(ctx, "voting-system-8b230", sa)
	})
	return firestoreClient, initErr
}
