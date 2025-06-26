package handlers

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"

	"go-voting-backend/eth"
)

type ExportedBlock struct {
	Number     uint64   `json:"number"`
	Hash       string   `json:"hash"`
	ParentHash string   `json:"parentHash"`
	Timestamp  uint64   `json:"timestamp"`
	TxCount    int      `json:"txCount"`
	Txs        []string `json:"txs"`
}

func BlockListHandler(w http.ResponseWriter, r *http.Request) {
	client, err := eth.GetClient()
	if err != nil {
		http.Error(w, "Eth client not initialized: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取最新区块号
	latest, err := client.Client.BlockNumber(context.Background())
	if err != nil {
		http.Error(w, "Failed to get latest block: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 遍历从创世块（区块号 0）到最新区块
	var blocks []ExportedBlock
	for i := int64(0); i <= int64(latest); i++ {
		block, err := client.Client.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			continue // 忽略失败的块
		}

		var txs []string
		for _, tx := range block.Transactions() {
			txs = append(txs, tx.Hash().Hex())
		}

		blocks = append(blocks, ExportedBlock{
			Number:     block.NumberU64(),
			Hash:       block.Hash().Hex(),
			ParentHash: block.ParentHash().Hex(),
			Timestamp:  block.Time(),
			TxCount:    len(block.Transactions()),
			Txs:        txs,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blocks)
}
