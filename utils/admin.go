// utils/admin.go
package utils

import (
	"context"
	"errors"
)

func GetAdminPassword() (string, error) {
	client, err := GetFirestoreClient()
	if err != nil {
		return "", err
	}

	doc, err := client.Collection("admins").Doc("admin1").Get(context.Background())
	if err != nil {
		return "", err
	}

	password, ok := doc.Data()["password"].(string)
	if !ok {
		return "", errors.New("invalid password format")
	}

	return password, nil
}
