package utils

import (
	"context"
	"errors"
	"fmt"
	"net"

	"cloud.google.com/go/firestore"
)

func CheckVoteEligibility(ic string, ip net.IP, client *firestore.Client) error {
	// 检查是否已投票（用 IC 作为主键）
	doc, err := client.Collection("votes").Doc(ic).Get(context.Background())
	if err == nil && doc.Exists() {
		return errors.New("🚫 You have already voted")
	}

	// 检查是否有相同 IP 的记录
	ipStr := ip.String()
	iter := client.Collection("votes").Where("ip", "==", ipStr).Documents(context.Background())
	defer iter.Stop()

	// 用 Next 判断是否存在至少一条记录
	if _, err := iter.Next(); err == nil {
		return fmt.Errorf("🚫 This IP %s has already been used to vote", ipStr)
	}

	return nil
}
