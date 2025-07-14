package utils

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	FirestoreClient *firestore.Client
	initOnce        sync.Once
	initErr         error
)

func InitFirestore() error {
	initOnce.Do(func() {
		ctx := context.Background()
		// ✅ 直接从 Secret File 路径加载
		sa := option.WithCredentialsFile("/etc/secrets/firebase-service-account.json")

		projectID := os.Getenv("GOOGLE_PROJECT_ID")
		FirestoreClient, initErr = firestore.NewClient(ctx, projectID, sa)
	})
	return initErr
}

func IsICRegistered(ic string) (bool, error) {
	doc, err := FirestoreClient.Collection("users").Doc(ic).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}
	return doc.Exists(), nil
}
