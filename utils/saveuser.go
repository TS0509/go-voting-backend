package utils

import (
	"context"

	"cloud.google.com/go/firestore"
)

func SaveUser(user User) error {
	data := map[string]interface{}{
		"ic":           user.IC,
		"privateKey":   user.PrivateKey,
		"address":      user.Address,
		"faceImageUrl": user.FaceImage, // 🔧 字段名必须和前端一致
		"hasVoted":     user.HasVoted,
		"lastIP":       user.LastIP,
	}
	_, err := FirestoreClient.Collection("users").Doc(user.IC).Set(
		context.Background(),
		data,
		firestore.MergeAll, // ✅ 合并模式，OK!
	)
	return err
}
