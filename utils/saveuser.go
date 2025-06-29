package utils

import (
	"context"
)

// 直接写入用户
func SaveUser(user User) error {
	client, err := GetFirestoreClient()
	if err != nil {
		return err
	}
	_, err = client.Collection("users").Doc(user.IC).Set(context.Background(), user)
	return err
}
