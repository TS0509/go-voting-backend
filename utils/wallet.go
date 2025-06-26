package utils

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func GenerateWallet() (*Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return &Wallet{Address: address, PrivateKey: "0x" + privateKeyHex}, nil
}
