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
	// 加载 .env
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

	// 加载私钥
	systemPriv, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid system private key: %v", err)
	}
	fromAddress := crypto.PubkeyToAddress(systemPriv.PublicKey)

	// 连接区块链客户端
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %v", err)
	}
	defer client.Close()

	// 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	// 获取 ChainID 和建议 gasPrice
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chainID: %v", err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}

	// ✅ 固定转账 0.0005 ETH（单位为 wei）
	value := new(big.Int)
	value.SetString("500000000000000", 10) // 0.0005 ETH

	txGasLimit := uint64(21000)

	fmt.Println("🔍 目标地址:", toAddress)
	fmt.Println("🔍 发起方地址:", fromAddress.Hex())
	fmt.Println("🔍 gasPrice:", gasPrice.String())
	fmt.Println("🔍 发送金额 (wei):", value.String())

	// 创建交易
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(toAddress),
		value,
		txGasLimit,
		gasPrice,
		nil,
	)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), systemPriv)
	if err != nil {
		return fmt.Errorf("signTx failed: %v", err)
	}

	// 广播交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("send transaction failed: %v", err)
	}

	log.Println("✅ 自动打币已发送，哈希:", signedTx.Hash().Hex())

	// 等待确认
	for {
		receipt, _ := client.TransactionReceipt(context.Background(), signedTx.Hash())
		if receipt != nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				log.Println("✅ 充值区块确认成功:", receipt.BlockNumber)
				break
			} else {
				return fmt.Errorf("❌ 交易上链但状态失败")
			}
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}
