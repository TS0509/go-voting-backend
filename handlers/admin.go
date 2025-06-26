package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"go-voting-backend/eth"
)

type AdminRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func AddCandidateHandler(w http.ResponseWriter, r *http.Request) {
	var req AdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Password != os.Getenv("ADMIN_PASSWORD") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	json.NewEncoder(w).Encode(map[string]string{
		"txHash": tx.Hash().Hex(),
	})
}

func StartVotingHandler(w http.ResponseWriter, r *http.Request) {
	var req AdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Password != os.Getenv("ADMIN_PASSWORD") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
