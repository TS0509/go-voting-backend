package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go-voting-backend/utils"

	"cloud.google.com/go/firestore"
)

type RegisterRequest struct {
	IC string `json:"ic"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req RegisterRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	wallet, err := utils.GenerateWallet()
	if err != nil {
		http.Error(w, "failed to generate wallet", http.StatusInternalServerError)
		return
	}
	client, err := firestore.NewClient(ctx, "voting-system-8b230")
	if err != nil {
		http.Error(w, "firestore error", http.StatusInternalServerError)
		return
	}
	defer client.Close()
	_, err = client.Collection("users").Doc(req.IC).Set(ctx, map[string]interface{}{
		"address":     wallet.Address,
		"private_key": wallet.PrivateKey,
	})
	if err != nil {
		http.Error(w, "save failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(wallet)
}
