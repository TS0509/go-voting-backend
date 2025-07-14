package handlers

import (
	"encoding/json"
	"net/http"

	"go-voting-backend/contract" // 假设这里是你生成的合约绑定
	"go-voting-backend/eth"
)

type VoteLog struct {
	Voter        string `json:"voter"`
	CandidateIdx uint64 `json:"candidateIndex"`
	TxHash       string `json:"txHash"`
	BlockNumber  uint64 `json:"blockNumber"`
}

func VoteLogHandler(w http.ResponseWriter, r *http.Request) {
	client, err := eth.GetClient()
	if err != nil {
		http.Error(w, "eth client error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	contractAddr := client.ContractAddress

	votingContract, err := contract.NewVoting(contractAddr, client.Client)
	if err != nil {
		http.Error(w, "contract bind error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logs, err := votingContract.FilterVoted(nil, nil, nil) // 过滤所有事件
	if err != nil {
		http.Error(w, "event filter error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var results []VoteLog
	for logs.Next() {
		event := logs.Event
		results = append(results, VoteLog{
			Voter:        event.Voter.Hex(),
			CandidateIdx: event.CandidateIndex.Uint64(),
			TxHash:       event.Raw.TxHash.Hex(),
			BlockNumber:  event.Raw.BlockNumber,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
