package utils

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type User struct {
	IC         string `firestore:"ic"`
	PrivateKey string `firestore:"privateKey"`
	HasVoted   bool   `firestore:"hasVoted"`
	LastIP     string `firestore:"lastIP"`
}

// 根据 IC 获取用户
func GetUserByIC(ic string) (*User, error) {
	client, err := GetFirestoreClient()
	if err != nil {
		return nil, err
	}

	doc, err := client.Collection("users").Doc(ic).Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("IC not found: %s", ic)
	}

	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// 更新为已投票状态
func MarkUserVoted(ic string) error {
	client, err := GetFirestoreClient()
	if err != nil {
		return err
	}

	_, err = client.Collection("users").Doc(ic).Update(context.Background(), []firestore.Update{
		{Path: "hasVoted", Value: true},
		{Path: "lastIP", Value: "127.0.0.1"}, // 🛠 这里你可以替换为真实 IP
	})

	return err
}
