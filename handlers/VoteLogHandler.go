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

	// ✅ 获取最新区块
	latestHeader, err := client.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("❌ 获取最新区块头失败:", err)
		http.Error(w, "get latest block error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	latestBlock := latestHeader.Number.Uint64()

	// ✅ 定义回溯区块数和每次查询跨度
	const contractDeployedAt uint64 = 8765000 // ⬅️ 这里换成你查到的部署区块号
	const step uint64 = 500
	startBlock := contractDeployedAt

	log.Printf("🔍 正在分段读取投票事件，起始区块 #%d -> 最新区块 #%d\n", startBlock, latestBlock)

	var results []VoteLog

	// ✅ 分段查询 logs
	for from := startBlock; from <= latestBlock; from += step {
		to := from + step - 1
		if to > latestBlock {
			to = latestBlock
		}

		opts := &bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: context.Background(),
		}

		iter, err := votingContract.FilterVoted(opts, nil, nil)
		if err != nil {
			log.Printf("❌ 查询区块 [%d ~ %d] 失败: %v", from, to, err)
			continue // 跳过失败的区段
		}

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
			log.Printf("⚠️ 事件迭代器错误（区块 %d ~ %d）: %v", from, to, iter.Error())
		}
		iter.Close()
	}

	log.Printf("📤 共返回 %d 条投票记录\n", len(results))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
