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

	latest, err := client.Client.BlockNumber(context.Background())
	if err != nil {
		http.Error(w, "Failed to get latest block: "+err.Error(), http.StatusInternalServerError)
		return
	}

	const MaxBlocks = 20
	start := int64(0)
	if int64(latest) > MaxBlocks {
		start = int64(latest) - MaxBlocks + 1
	}

	var blocks []ExportedBlock
	for i := start; i <= int64(latest); i++ {
		block, err := client.Client.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			continue
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
