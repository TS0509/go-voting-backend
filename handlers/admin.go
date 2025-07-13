package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go-voting-backend/eth"
	"go-voting-backend/middleware"
	"go-voting-backend/utils"
)

type AddCandidateRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func AddCandidateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📩 AddCandidateHandler triggered")

	role := r.Context().Value(middleware.RoleKey)
	if role != "admin" {
		http.Error(w, "Forbidden: not admin", http.StatusForbidden)
		return
	}

	var req AddCandidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Avatar == "" {
		http.Error(w, "Name and avatar are required", http.StatusBadRequest)
		return
	}

	client, err := eth.GetClient()
	if err != nil {
		http.Error(w, "Eth client not initialized", http.StatusInternalServerError)
		return
	}

	auth, err := client.GetAuth()
	if err != nil {
		http.Error(w, "Auth error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := client.Contract.AddCandidate(auth, req.Name)
	if err != nil {
		http.Error(w, "Add candidate failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ 保存头像等候选人数据到 Firestore
	if req.Avatar != "" {
		_, _, err := utils.FirestoreClient.Collection("candidates").Add(context.Background(), map[string]interface{}{
			"name":   req.Name,
			"avatar": req.Avatar,
			"txHash": tx.Hash().Hex(),
		})
		if err != nil {
			http.Error(w, "Firestore save failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"txHash": tx.Hash().Hex(),
	})
}

func StartVotingHandler(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey)
	if role != "admin" {
		http.Error(w, "Forbidden: not admin", http.StatusForbidden)
		return
	}

	client, err := eth.GetClient()
	if err != nil {
		http.Error(w, "Eth client not initialized", http.StatusInternalServerError)
		return
	}

	auth, err := client.GetAuth()
	if err != nil {
		http.Error(w, "Auth error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := client.Contract.StartVoting(auth)
	if err != nil {
		http.Error(w, "Start voting failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"txHash": tx.Hash().Hex(),
	})
}

func StopVotingHandler(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey)
	if role != "admin" {
		http.Error(w, "Forbidden: not admin", http.StatusForbidden)
		return
	}

	client, err := eth.GetClient()
	if err != nil {
		http.Error(w, "Eth client not initialized", http.StatusInternalServerError)
		return
	}

	auth, err := client.GetAuth()
	if err != nil {
		http.Error(w, "Auth error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := client.Contract.EndVoting(auth)
	if err != nil {
		http.Error(w, "End voting failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"txHash": tx.Hash().Hex(),
	})
}
