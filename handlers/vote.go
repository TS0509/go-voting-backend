package handlers

import (
	"encoding/json"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"

	"go-voting-backend/eth"
	"go-voting-backend/utils"

	"github.com/ethereum/go-ethereum/common"
)

type VoteRequest struct {
	CandidateIndex int    `json:"candidate"`
	IC             string `json:"ic"`
}

func getRealIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "127.0.0.1"
	}
	return host
}

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("➡️ 收到 /vote 请求")

	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("❌ 请求解析失败:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	log.Println("🔍 投票请求 - IC:", req.IC, "候选人索引:", req.CandidateIndex)

	user, err := utils.GetUserByIC(req.IC)
	if err != nil || user == nil {
		log.Println("❌ 用户不存在或查询失败:", err)
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}
	log.Printf("✅ Firestore 用户获取成功: %+v", user)

	// ✅ 加这一段 —— 防止重复投票
	if user.HasVoted {
		log.Println("❌ 用户已投过票，拒绝重复投票")
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}

	rpc := os.Getenv("RPC_URL")
	contract := os.Getenv("CONTRACT_ADDRESS")
	if rpc == "" || contract == "" {
		log.Println("❌ 环境变量 RPC_URL 或 CONTRACT_ADDRESS 未配置")
		http.Error(w, "env config missing", http.StatusInternalServerError)
		return
	}

	client, err := eth.NewEthClient(rpc, common.HexToAddress(contract), user.PrivateKey)
	if err != nil {
		log.Println("❌ 初始化 EthClient 失败:", err)
		http.Error(w, "init eth client failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ipStr := getRealIP(r)
	ip := net.ParseIP(ipStr)
	if ip == nil {
		ip = net.IPv4(127, 0, 0, 1)
	}
	log.Println("🔍 用户 IP:", ipStr)

	// ✅ 额外限制：IP + IC 重复
	if err := utils.CheckVoteEligibility(req.IC, ip, utils.FirestoreClient); err != nil {
		log.Println("❌ 投票资格检查失败:", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	started, err := client.Contract.HasStarted(nil)
	if err != nil {
		log.Println("❌ 合约投票状态检查失败:", err)
		http.Error(w, "check voting failed", http.StatusInternalServerError)
		return
	}
	if !started {
		log.Println("⚠️ 投票尚未开始")
		http.Error(w, "voting not started", http.StatusForbidden)
		return
	}
	log.Println("✅ 投票已启动")

	auth, err := client.GetAuth()
	if err != nil {
		log.Println("❌ 获取交易身份失败:", err)
		http.Error(w, "auth error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := client.Contract.Vote(auth, big.NewInt(int64(req.CandidateIndex)))
	if err != nil {
		log.Println("❌ vote failed:", err)
		http.Error(w, "vote failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("✅ 投票成功，交易哈希:", tx.Hash().Hex())

	markErr := utils.MarkUserVoted(req.IC, ipStr)
	if markErr != nil {
		log.Println("⚠️ 已投票记录失败:", markErr)
	}

	resp := map[string]string{
		"txHash": tx.Hash().Hex(),
	}
	if markErr != nil {
		resp["warning"] = "Vote succeeded, but failed to mark user as voted"
	}
	json.NewEncoder(w).Encode(resp)
}
