package utils

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func FundWallet(toAddress string) error {
	// .env
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env: %v", err)
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return fmt.Errorf("missing PRIVATE_KEY in .env")
	}
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		return fmt.Errorf("missing RPC_URL in .env")
	}

	systemPriv, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid system private key: %v", err)
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %v", err)
	}
	defer client.Close()

	fromAddress := crypto.PubkeyToAddress(systemPriv.PublicKey)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chainID: %v", err)
	}

	value := big.NewInt(1e16) // 0.01 ETH
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}

	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(toAddress),
		value,
		gasLimit,
		gasPrice,
		nil,
	)

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), systemPriv)
	if err != nil {
		return fmt.Errorf("signTx failed: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("send transaction failed: %v", err)
	}

	log.Println("自动打币已发送，哈希:", signedTx.Hash().Hex())

	// 等确认
	for {
		receipt, _ := client.TransactionReceipt(context.Background(), signedTx.Hash())
		if receipt != nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				log.Println("充值区块确认:", receipt.BlockNumber)
				break
			} else {
				return fmt.Errorf("交易上链但状态失败")
			}
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}
