package handlers

import (
	"encoding/json"
	"go-voting-backend/utils"
	"net/http"
)

type RegisterRequest struct {
	IC string `json:"ic"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// 创建新钱包
	wallet, err := utils.GenerateWallet()
	if err != nil {
		http.Error(w, "failed to generate wallet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 自动打币
	if err := utils.FundWallet(wallet.Address); err != nil {
		http.Error(w, "funding failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 写 Firestore
	user := utils.User{
		IC:         req.IC,
		PrivateKey: wallet.PrivateKey,
		Address:    wallet.Address,
		FaceImage:  "",
		HasVoted:   false,
		LastIP:     "",
	}

	if err := utils.SaveUser(user); err != nil {
		http.Error(w, "failed to save user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"address": wallet.Address,
	})
}
