package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go-voting-backend/eth"
	"go-voting-backend/utils"
)

type CandidateInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Index  int    `json:"index"`
}

// GET /candidates
func GetCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	client, err := eth.GetClient()
	if err != nil {
		http.Error(w, "Eth client not initialized", http.StatusInternalServerError)
		return
	}

	candidates, err := client.Contract.GetCandidates(nil)
	if err != nil {
		http.Error(w, "Failed to get candidates from contract: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var results []CandidateInfo

	for i, c := range candidates {
		name := c.Name

		// 🔍 从 Firestore 查头像
		docsnap, err := utils.FirestoreClient.Collection("candidates").
			Where("name", "==", name).
			Limit(1).
			Documents(context.Background()).
			Next()

		avatar := ""
		if err == nil && docsnap != nil {
			data := docsnap.Data()
			if val, ok := data["avatar"].(string); ok {
				avatar = val
			}
		}

		results = append(results, CandidateInfo{
			Name:   name,
			Avatar: avatar,
			Index:  i,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
