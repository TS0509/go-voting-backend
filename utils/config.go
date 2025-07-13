package utils

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
)

var RPC_URL = os.Getenv("RPC_URL")

var ContractAddress = common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
