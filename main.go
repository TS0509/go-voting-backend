package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"go-voting-backend/eth"
	"go-voting-backend/handlers"
	"go-voting-backend/middleware"
	"go-voting-backend/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

// 🌐 支持环境变量配置的 CORS 中间件
func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// 是否匹配白名单中的 Origin
		for _, allowed := range allowedOrigins {
			if origin == strings.TrimSpace(allowed) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// ✅ 加载 .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env 文件未找到，尝试使用环境变量运行（适用于部署环境）")
	}

	// ✅ 初始化 Firestore
	if err := utils.InitFirestore(); err != nil {
		log.Fatal("❌ Firestore 初始化失败:", err)
	}

	// ✅ 初始化以太坊客户端
	rpcURL := os.Getenv("RPC_URL")
	privateKey := os.Getenv("PRIVATE_KEY")
	contractAddr := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	if err := eth.InitClient(rpcURL, contractAddr, privateKey); err != nil {
		log.Fatal("❌ InitClient failed:", err)
	}

	// ✅ 设置路由
	mux := http.NewServeMux()

	// ✅ 限流中间件（每 IP 每秒 5 次，突发容量 10）
	ipLimiter := middleware.NewIPLimiter(rate.Limit(5), 10)
	protected := ipLimiter.RateLimitMiddleware(mux)

	// ✅ 普通接口
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/votelog", handlers.VoteLogHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/vote", handlers.VoteHandler)
	mux.HandleFunc("/api/blocks", handlers.BlockListHandler)
	mux.Handle("/candidates", http.HandlerFunc(handlers.GetCandidatesHandler))

	// ✅ 管理员权限接口（需认证）
	mux.Handle("/admin/add-candidate", middleware.AuthMiddleware(http.HandlerFunc(handlers.AddCandidateHandler)))
	mux.Handle("/admin/start-voting", middleware.AuthMiddleware(http.HandlerFunc(handlers.StartVotingHandler)))
	mux.Handle("/admin/stop-voting", middleware.AuthMiddleware(http.HandlerFunc(handlers.StopVotingHandler)))

	// ✅ Token 校验接口
	mux.Handle("/auth/check", middleware.AuthMiddleware(http.HandlerFunc(handlers.AuthCheckHandler)))

	// ✅ 启动服务器
	if external := os.Getenv("RENDER_EXTERNAL_URL"); external != "" {
		log.Println("✅ Server deployed at:", external)
	} else {
		log.Println("✅ Server running at http://localhost:8080")
	}

	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(protected)))
}
