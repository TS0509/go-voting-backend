package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go-voting-backend/contract"
	"go-voting-backend/eth"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type VoteLog struct {
	Voter        string `json:"voter"`
	CandidateIdx uint64 `json:"candidateIndex"`
	TxHash       string `json:"txHash"`
	BlockNumber  uint64 `json:"blockNumber"`
}

func VoteLogHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📥 [VoteLogHandler] 接收到请求")

	client, err := eth.GetClient()
	if err != nil {
		log.Println("❌ 获取以太坊客户端失败:", err)
		http.Error(w, "eth client error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	contractAddr := client.ContractAddress

	votingContract, err := contract.NewVoting(contractAddr, client.Client)
	if err != nil {
		log.Println("❌ 绑定 Voting 合约失败:", err)
		http.Error(w, "contract bind error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ 获取当前最新区块头
	header, err := client.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("❌ 获取最新区块头失败:", err)
		http.Error(w, "get latest block error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ 向前回溯最多 N 个区块，避免全链遍历超时
	const blockLookback uint64 = 3000
	var fromBlock uint64 = 0
	if header.Number.Uint64() > blockLookback {
		fromBlock = header.Number.Uint64() - blockLookback
	}
	log.Printf("🔍 正在从区块 #%d 读取投票事件...\n", fromBlock)

	opts := &bind.FilterOpts{
		Start:   fromBlock,
		Context: context.Background(),
	}

	iter, err := votingContract.FilterVoted(opts, nil, nil)
	if err != nil {
		log.Println("❌ 读取投票事件失败:", err)
		http.Error(w, "event filter error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer iter.Close()

	var results []VoteLog

	for iter.Next() {
		event := iter.Event
		if event == nil || event.CandidateIndex == nil {
			log.Println("⚠️ 遇到无效事件，跳过")
			continue
		}
		log.Printf("✅ 捕获投票事件 - Voter: %s, Candidate: %d\n", event.Voter.Hex(), event.CandidateIndex.Uint64())

		results = append(results, VoteLog{
			Voter:        event.Voter.Hex(),
			CandidateIdx: event.CandidateIndex.Uint64(),
			TxHash:       event.Raw.TxHash.Hex(),
			BlockNumber:  event.Raw.BlockNumber,
		})
	}

	if iter.Error() != nil {
		log.Println("❌ 事件迭代器错误:", iter.Error())
		http.Error(w, "event iterate error: "+iter.Error().Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("📤 共返回 %d 条投票记录\n", len(results))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
