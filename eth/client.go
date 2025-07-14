package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"strings"
	"sync"

	"go-voting-backend/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	Client          *ethclient.Client
	Contract        *contract.Voting
	ContractAddress common.Address
	PrivKey         *ecdsa.PrivateKey
	ChainID         *big.Int
}

var (
	globalClient *EthClient
	once         sync.Once
)

// ✅ 用于后台初始化固定托管钱包（用于管理员操作）
func InitClient(rpcURL string, contractAddr common.Address, privateKey string) error {
	var err error
	once.Do(func() {
		client, errDial := ethclient.Dial(rpcURL)
		if errDial != nil {
			err = errDial
			return
		}

		priv, errKey := crypto.HexToECDSA(privateKey)
		if errKey != nil {
			err = errKey
			return
		}

		contractInstance, errContract := contract.NewVoting(contractAddr, client)
		if errContract != nil {
			err = errContract
			return
		}

		chainIDStr := os.Getenv("CHAIN_ID")
		chainID := big.NewInt(31337) // 默认
		if chainIDStr != "" {
			chainID.SetString(chainIDStr, 10)
		}

		globalClient = &EthClient{
			Client:          client,
			Contract:        contractInstance,
			ContractAddress: contractAddr,
			PrivKey:         priv,
			ChainID:         chainID,
		}
	})
	return err
}

func GetClient() (*EthClient, error) {
	if globalClient == nil {
		return nil, errors.New("❌ EthClient 未初始化")
	}
	return globalClient, nil
}

// ✅ 每个用户调用时用此函数，私钥是动态传入的
func NewEthClient(rpcURL string, contractAddr common.Address, privateKey string) (*EthClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	privateKey = strings.TrimPrefix(privateKey, "0x")

	priv, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	contractInstance, err := contract.NewVoting(contractAddr, client)
	if err != nil {
		return nil, err
	}

	chainIDStr := os.Getenv("CHAIN_ID")
	chainID := big.NewInt(31337)
	if chainIDStr != "" {
		chainID.SetString(chainIDStr, 10)
	}

	return &EthClient{
		Client:          client,
		Contract:        contractInstance,
		ContractAddress: contractAddr,
		PrivKey:         priv,
		ChainID:         chainID,
	}, nil
}

func (e *EthClient) GetAuth() (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(e.PrivKey, e.ChainID)
	if err != nil {
		return nil, err
	}

	auth.Context = context.Background()

	auth.GasLimit = uint64(100_000)
	return auth, nil
}
