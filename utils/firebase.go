package utils

import (
	"context"
	"errors"
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

// ✅ 从环境变量中读取服务账号 JSON 字符串 → 写入临时文件 → 用于认证
func InitFirestore() error {
	initOnce.Do(func() {
		ctx := context.Background()

		// 从环境变量获取 JSON
		credJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
		if credJSON == "" {
			initErr = errors.New("GOOGLE_APPLICATION_CREDENTIALS_JSON not set")
			return
		}

		// 写入临时文件
		tmpFile, err := os.CreateTemp("", "firebase-creds-*.json")
		if err != nil {
			initErr = err
			return
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(credJSON)); err != nil {
			initErr = err
			return
		}
		tmpFile.Close()

		sa := option.WithCredentialsFile(tmpFile.Name())
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
