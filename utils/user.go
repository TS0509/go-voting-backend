package utils

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type User struct {
	IC         string `firestore:"ic"`
	PrivateKey string `firestore:"privateKey"`
	Address    string `firestore:"address"`
	FaceImage  string `firestore:"faceImageUrl"`
	HasVoted   bool   `firestore:"hasVoted"`
	LastIP     string `firestore:"lastIP"`
}

// 根据 IC 获取用户
func GetUserByIC(ic string) (*User, error) {
	doc, err := FirestoreClient.Collection("users").Doc(ic).Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get user by IC (%s): %v", ic, err)
	}

	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %v", err)
	}

	return &user, nil
}

// 标记用户已投票，记录真实 IP
func MarkUserVoted(ic string, ip string) error {
	_, err := FirestoreClient.Collection("users").Doc(ic).Update(context.Background(), []firestore.Update{
		{Path: "hasVoted", Value: true},
		{Path: "lastIP", Value: ip},
	})
	if err != nil {
		return fmt.Errorf("failed to mark user voted: %v", err)
	}
	return nil
}
