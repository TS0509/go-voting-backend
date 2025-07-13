package utils

import (
	"context"
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
		sa := option.WithCredentialsFile("config/firebase-service-account.json")
		FirestoreClient, initErr = firestore.NewClient(ctx, "voting-system-8b230", sa)
	})
	return initErr
}

func IsICRegistered(ic string) (bool, error) {
	doc, err := FirestoreClient.Collection("users").Doc(ic).Get(context.Background())
	if err != nil {
		// 🔧 修复点：判断是否 NotFound
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}
	return doc.Exists(), nil
}
