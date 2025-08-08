package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

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

type PaginatedVoteLogs struct {
	Logs       []VoteLog `json:"logs"`
	TotalCount int       `json:"totalCount"`
	Page       int       `json:"page"`
	PageSize   int       `json:"pageSize"`
}

var contractDeployedAt uint64
var step uint64

func init() {
	if val := os.Getenv("CONTRACT_DEPLOYED_AT"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 64); err == nil {
			contractDeployedAt = num
		} else {
			log.Fatalf("❌ 无效的 CONTRACT_DEPLOYED_AT 值: %v", err)
		}
	} else {
		// 没设置就给个默认值（比如本地测试用）
		contractDeployedAt = 8939266
		log.Println("⚠️ 未设置 CONTRACT_DEPLOYED_AT，使用默认值 8939266")
	}

	if val := os.Getenv("STEP"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 64); err == nil {
			step = num
		} else {
			log.Fatalf("❌ 无效的 STEP 值: %v", err)
		}
	} else {
		step = 500
	}
}

func VoteLogHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📥 [VoteLogHandler] 接收到请求")

	// 连接客户端
	client, err := eth.GetClient()
	if err != nil {
		log.Println("❌ 获取以太坊客户端失败:", err)
		http.Error(w, "eth client error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	contractAddr := client.ContractAddress

	// 绑定合约
	votingContract, err := contract.NewVoting(contractAddr, client.Client)
	if err != nil {
		log.Println("❌ 绑定 Voting 合约失败:", err)
		http.Error(w, "contract bind error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取最新区块
	latestHeader, err := client.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("❌ 获取最新区块头失败:", err)
		http.Error(w, "get latest block error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	latestBlock := latestHeader.Number.Uint64()
	log.Printf("🧱 最新区块高度: %d\n", latestBlock)

	// 分页参数处理
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if pageSize < 1 {
		pageSize = 10
	}
	log.Printf("📄 请求分页参数: page=%d, size=%d\n", page, pageSize)

	var allLogs []VoteLog
	log.Printf("🔍 正在分段读取投票事件，起始区块 #%d -> 最新区块 #%d\n", contractDeployedAt, latestBlock)

	// 批量查询 logs
	for from := contractDeployedAt; from <= latestBlock; from += step {
		to := from + step - 1
		if to > latestBlock {
			to = latestBlock
		}
		log.Printf("📦 查询区块 [%d - %d]", from, to)

		opts := &bind.FilterOpts{Start: from, End: &to, Context: context.Background()}
		iter, err := votingContract.FilterVoted(opts, nil, nil)
		if err != nil {
			log.Printf("⚠️ 查询失败 [%d - %d]: %v", from, to, err)
			continue
		}
		for iter.Next() {
			event := iter.Event
			if event == nil || event.CandidateIndex == nil {
				log.Println("⚠️ 遇到无效事件，跳过")
				continue
			}
			log.Printf("✅ 捕获投票事件 - Voter: %s, Candidate: %d, Block: %d",
				event.Voter.Hex(), event.CandidateIndex.Uint64(), event.Raw.BlockNumber)

			allLogs = append(allLogs, VoteLog{
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

	total := len(allLogs)
	log.Printf("📊 共捕获 %d 条投票记录", total)

	// 分页截取
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginated := PaginatedVoteLogs{
		Logs:       allLogs[start:end],
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}

	log.Printf("📤 正在返回分页数据 [%d ~ %d)", start, end)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginated)
}
