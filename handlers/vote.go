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

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := utils.GetUserByIC(req.IC)
	if err != nil || user == nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	rpc := os.Getenv("RPC_URL")
	contract := os.Getenv("CONTRACT_ADDRESS")
	if rpc == "" || contract == "" {
		http.Error(w, "env config missing", http.StatusInternalServerError)
		return
	}

	client, err := eth.NewEthClient(rpc, common.HexToAddress(contract), user.PrivateKey)
	if err != nil {
		http.Error(w, "init eth client failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 解析ip
	ipStr := strings.Split(r.RemoteAddr, ":")[0]
	ip := net.ParseIP(ipStr)
	if ip == nil {
		ip = net.IPv4(127, 0, 0, 1)
	}

	// eligibility
	dbClient, err := utils.GetFirestoreClient()
	if err != nil {
		http.Error(w, "firestore error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := utils.CheckVoteEligibility(req.IC, ip, dbClient); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	started, err := client.Contract.HasStarted(nil)
	if err != nil {
		http.Error(w, "check voting failed", http.StatusInternalServerError)
		return
	}
	if !started {
		http.Error(w, "voting not started", http.StatusForbidden)
		return
	}

	auth, err := client.GetAuth()
	if err != nil {
		http.Error(w, "auth error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tx, err := client.Contract.Vote(auth, big.NewInt(int64(req.CandidateIndex)))
	if err != nil {
		http.Error(w, "vote failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := utils.MarkUserVoted(req.IC); err != nil {
		http.Error(w, "vote success but mark error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"txHash": tx.Hash().Hex(),
	})
}
