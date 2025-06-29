package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"

	"go-voting-backend/eth"
	"go-voting-backend/handlers"
)

// 🌐 通用 CORS 中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // 可改为特定前端网址
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Failed to load .env")
	}

	rpcURL := os.Getenv("RPC_URL")
	privateKey := os.Getenv("PRIVATE_KEY")
	contractAddr := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	err = eth.InitClient(rpcURL, contractAddr, privateKey)
	if err != nil {
		log.Fatal("❌ InitClient failed:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/vote", handlers.VoteHandler)
	mux.HandleFunc("/api/blocks", handlers.BlockListHandler)
	mux.HandleFunc("/admin/add-candidate", handlers.AddCandidateHandler)
	mux.HandleFunc("/admin/start-voting", handlers.StartVotingHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)

	log.Println("✅ Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}
