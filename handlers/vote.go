package handlers

import (
	"encoding/json"
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

func respondError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// 获取用户
	user, err := utils.GetUserByIC(req.IC)
	if err != nil {
		respondError(w, http.StatusBadRequest, "User not found: "+err.Error())
		return
	}

	// 获取私钥钱包对应客户端
	rpc := os.Getenv("RPC_URL")
	contractAddr := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	client, err := eth.NewEthClient(rpc, contractAddr, user.PrivateKey)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to init eth client: "+err.Error())
		return
	}

	// IP 地址解析
	ipStr := strings.Split(r.RemoteAddr, ":")[0]
	ip := net.ParseIP(ipStr)
	if ip == nil {
		ip = net.IPv4(127, 0, 0, 1)
	}

	// 检查投票资格
	dbClient, err := utils.GetFirestoreClient()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Firestore error: "+err.Error())
		return
	}
	if err := utils.CheckVoteEligibility(req.IC, ip, dbClient); err != nil {
		respondError(w, http.StatusForbidden, err.Error())
		return
	}

	// 检查投票是否开始
	started, err := client.Contract.HasStarted(nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Check voting status error: "+err.Error())
		return
	}
	if !started {
		respondError(w, http.StatusForbidden, "Voting not started yet")
		return
	}

	// 创建授权并发起投票
	auth, err := client.GetAuth()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Auth error: "+err.Error())
		return
	}

	tx, err := client.Contract.Vote(auth, big.NewInt(int64(req.CandidateIndex)))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Vote failed: "+err.Error())
		return
	}

	// 标记已投票
	if err := utils.MarkUserVoted(req.IC); err != nil {
		respondError(w, http.StatusInternalServerError, "Vote succeeded but marking failed: "+err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"txHash": tx.Hash().Hex(),
	})
}
